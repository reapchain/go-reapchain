// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/consensus/podc/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	elog "github.com/ethereum/go-ethereum/log"
)

var testLogger = elog.New()

type testSystemBackend struct {
	id  uint64
	sys *testSystem

	engine Engine
	peers  podc.ValidatorSet
	events *event.TypeMux

	commitMsgs     []podc.Proposal
	committedSeals [][]byte
	sentMsgs       [][]byte // store the message when Send is called by core

	address common.Address
	db      ethdb.Database
}

// ==============================================
//
// define the functions that needs to be provided for Istanbul.

func (self *testSystemBackend) Address() common.Address {
	return self.address
}

// Peers returns all connected peers
func (self *testSystemBackend) Validators(proposal podc.Proposal) podc.ValidatorSet {
	return self.peers
}

func (self *testSystemBackend) EventMux() *event.TypeMux {
	return self.events
}

func (self *testSystemBackend) Send(message []byte, target common.Address) error {
	testLogger.Info("enqueuing a message...", "address", self.Address())
	self.sentMsgs = append(self.sentMsgs, message)
	self.sys.queuedMessage <- podc.MessageEvent{
		Payload: message,
	}
	return nil
}

func (self *testSystemBackend) Broadcast(valSet podc.ValidatorSet, message []byte) error {
	testLogger.Info("enqueuing a message...", "address", self.Address())
	self.sentMsgs = append(self.sentMsgs, message)
	self.sys.queuedMessage <- podc.MessageEvent{
		Payload: message,
	}
	return nil
}


// Multicast sends a message to specific targets
func (self *testSystemBackend) Multicast( payload []byte, targets []common.Address ) error {
	testLogger.Info("enqueuing a message...", "address", self.Address())
	self.sentMsgs = append(self.sentMsgs, payload)
	self.sys.queuedMessage <- podc.MessageEvent{
		Payload: payload,
	}
	return nil
}
//end


func (self *testSystemBackend) NextRound() error {
	testLogger.Warn("nothing to happen")
	return nil
}

func (self *testSystemBackend) Commit(proposal podc.Proposal, seals []byte) error {
	testLogger.Info("commit message", "address", self.Address())
	self.commitMsgs = append(self.commitMsgs, proposal)
	self.committedSeals = append(self.committedSeals, seals)

	// fake new head events
	go self.events.Post(podc.FinalCommittedEvent{
		Proposal: proposal,
	})
	return nil
}

func (self *testSystemBackend) Verify(proposal podc.Proposal) error {
	return nil
}

func (self *testSystemBackend) Sign(data []byte) ([]byte, error) {
	testLogger.Warn("not sign any data")
	return data, nil
}

func (self *testSystemBackend) CheckSignature([]byte, common.Address, []byte) error {
	return nil
}

func (self *testSystemBackend) CheckValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return common.Address{}, nil
}

func (self *testSystemBackend) Hash(b interface{}) common.Hash {
	return common.StringToHash("Test")
}

func (self *testSystemBackend) NewRequest(request podc.Proposal) {
	go self.events.Post(podc.RequestEvent{
		Proposal: request,
	})
}

// define the struct that need to be provided for integration tests.

type testSystem struct {
	backends []*testSystemBackend

	queuedMessage chan podc.MessageEvent
	quit          chan struct{}
}

func newTestSystem(n uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)
	return &testSystem{
		backends: make([]*testSystemBackend, n),

		queuedMessage: make(chan podc.MessageEvent),
		quit:          make(chan struct{}),
	}
}

func generateValidators(n int) []common.Address {
	vals := make([]common.Address, 0)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		vals = append(vals, crypto.PubkeyToAddress(privateKey.PublicKey))
	}
	return vals
}

func newTestValidatorSet(n int) podc.ValidatorSet {
	return validator.NewSet(generateValidators(n), podc.RoundRobin)
}

// FIXME: int64 is needed for N and F
func NewTestSystemWithBackend(n, f uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)

	addrs := generateValidators(int(n))
	sys := newTestSystem(n)
	config := podc.DefaultConfig

	for i := uint64(0); i < n; i++ {
		vset := validator.NewSet(addrs, podc.RoundRobin)
		backend := sys.NewBackend(i)
		backend.peers = vset
		backend.address = vset.GetByIndex(i).Address()

		core := New(backend, config).(*core)
		core.state = StateAcceptRequest
		core.lastProposer = common.Address{}
		core.current = newRoundState(&podc.View{
			Round:    big.NewInt(0),
			Sequence: big.NewInt(1),
		}, vset)
		core.logger = testLogger
		core.validateFn = backend.CheckValidatorSignature

		backend.engine = core
	}

	return sys
}

// listen will consume messages from queue and deliver a message to core
func (t *testSystem) listen() {
	for {
		select {
		case <-t.quit:
			return
		case queuedMessage := <-t.queuedMessage:
			testLogger.Info("consuming a queue message...")
			for _, backend := range t.backends {
				go backend.EventMux().Post(queuedMessage)
			}
		}
	}
}

// Run will start system components based on given flag, and returns a closer
// function that caller can control lifecycle
//
// Given a true for core if you want to initialize core engine.
func (t *testSystem) Run(core bool) func() {

	var qman []*discover.Node
	for _, b := range t.backends {
		if core {
			b.engine.Start(common.Big0, common.Address{}, nil, qman) // start PoDC core
		}
	}

	go t.listen()
	closer := func() { t.stop(core) }
	return closer
}

func (t *testSystem) stop(core bool) {
	close(t.quit)

	for _, b := range t.backends {
		if core {
			b.engine.Stop()
		}
	}
}

func (t *testSystem) NewBackend(id uint64) *testSystemBackend {
	// assume always success
	ethDB, _ := ethdb.NewMemDatabase()
	backend := &testSystemBackend{
		id:     id,
		sys:    t,
		events: new(event.TypeMux),
		db:     ethDB,
	}

	t.backends[id] = backend
	return backend
}

// helper functions.

func getPublicKeyAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
