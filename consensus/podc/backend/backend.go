
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

package backend

import (
"crypto/ecdsa"
"sync"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/consensus"
"github.com/ethereum/go-ethereum/consensus/podc"
 podcCore "github.com/ethereum/go-ethereum/consensus/podc/core"
"github.com/ethereum/go-ethereum/consensus/podc/validator"
"github.com/ethereum/go-ethereum/core"
"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/crypto"
"github.com/ethereum/go-ethereum/ethdb"
"github.com/ethereum/go-ethereum/event"
"github.com/ethereum/go-ethereum/log"
lru "github.com/hashicorp/golang-lru"
)

func New(config *podc.Config, eventMux *event.TypeMux, privateKey *ecdsa.PrivateKey, db ethdb.Database) consensus.PoDC {
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	backend := &simpleBackend{
		config:           config,
		eventMux:         eventMux,
		//istanbulEventMux: new(event.TypeMux),
		podcEventMux: new(event.TypeMux),
		privateKey:       privateKey,
		address:          crypto.PubkeyToAddress(privateKey.PublicKey),
		logger:           log.New("backend", "simple"),
		db:               db,
		commitCh:         make(chan *types.Block, 1),
		recents:          recents,
		candidates:       make(map[common.Address]bool),
	}
	return backend
}

// ----------------------------------------------------------------------------

type simpleBackend struct {
	config           *podc.Config
	eventMux         *event.TypeMux
	podcEventMux     *event.TypeMux
	privateKey       *ecdsa.PrivateKey
	address          common.Address

	core             podcCore.Engine
	logger           log.Logger
	quitSync         chan struct{}
	db               ethdb.Database
	timeout          uint64
	chain            consensus.ChainReader
	inserter         func(block *types.Block) error

	// the channels for istanbul engine notifications
	commitCh          chan *types.Block
	proposedBlockHash common.Hash
	sealMu            sync.Mutex

	// Current list of candidates we are pushing
	candidates map[common.Address]bool
	// Protects the signer fields
	candidatesLock sync.RWMutex
	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache
}

// Address implements istanbul.Backend.Address
func (sb *simpleBackend) Address() common.Address {
	return sb.address
}


//Validators implements PoDC.Backend.Validators
func (sb *simpleBackend) Validators(proposal podc.Proposal) podc.ValidatorSet {
	snap, err := sb.snapshot(sb.chain, proposal.Number().Uint64(), proposal.Hash(), nil)
	if err != nil {
		return validator.NewSet(nil, sb.config.ProposerPolicy)  // 지금은 라운드로빈으로, 동작 예정, 
	}
	return snap.ValSet
}

func (sb *simpleBackend) Send(payload []byte, target common.Address) error {
	go sb.eventMux.Post(podc.ConsensusDataEvent{
		Target: target,
		Data:   payload,
	})
	return nil
}


//

// Broadcast implements podc.Backend.Send
func (sb *simpleBackend) Broadcast(valSet podc.ValidatorSet, payload []byte) error {

	// 모든 Validator 리스트에 다 보냄.
	for _, val := range valSet.List() {
		if val.Address() == sb.Address() {
			// send to self
			msg := podc.MessageEvent{
				Payload: payload,
			}
			go sb.podcEventMux.Post(msg)

		} else {
			// send to other peers
			sb.Send(payload, val.Address())

		}
	}
	return nil
}

func (sb *simpleBackend) Multicast(payload []byte, targets []common.Address) error {
	for _, addr := range targets {
		if addr == sb.Address() {
			msg := podc.MessageEvent{
				Payload: payload,
			}
			go sb.podcEventMux.Post(msg)
		} else {
			sb.Send(payload, addr)
		}
	}
	return nil
}

// Commit implements podc.Backend.Commit
func (sb *simpleBackend) Commit(proposal podc.Proposal, seals []byte) error {
	// Check if the proposal is a valid block
	block := &types.Block{}
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return errInvalidProposal
	}

	h := block.Header()
	// Append seals into extra-data
	err := writeCommittedSeals(h, seals)
	if err != nil {
		return err
	}
	// update block's header
	block = block.WithSeal(h)

	sb.logger.Info("Committed", "address", sb.Address(), "hash", proposal.Hash(), "number", proposal.Number().Uint64())
	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == block.Hash() {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- block
		// TODO: how do we check the block is inserted correctly?
		return nil
	}

	return sb.inserter(block)
}

// NextRound will broadcast ChainHeadEvent to trigger next seal()

func (sb *simpleBackend) NextRound() error {
	header := sb.chain.CurrentHeader()
	sb.logger.Debug("NextRound", "address", sb.Address(), "current_hash", header.Hash(), "current_number", header.Number)
	go sb.eventMux.Post(core.ChainHeadEvent{})
	return nil
}

// EventMux implements PoDC.Backend.EventMux
func (sb *simpleBackend) EventMux() *event.TypeMux {
	return sb.podcEventMux
}

// Verify implements podc.Backend.Verify
func (sb *simpleBackend) Verify(proposal podc.Proposal) error {
	// Check if the proposal is a valid block
	block := &types.Block{}
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return errInvalidProposal
	}
	// verify the header of proposed block
	err := sb.VerifyHeader(sb.chain, block.Header(), false)
	// Ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err != nil && err != errEmptyCommittedSeals {
		return err
	}
	return nil
}

// Sign implements podc.Backend.Sign
func (sb *simpleBackend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256([]byte(data))
	return crypto.Sign(hashData, sb.privateKey)
}

func (sb *simpleBackend) CheckSignature(data []byte, address common.Address, sig []byte) error {
	signer, err := podc.GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return err
	}
	// Compare derived addresses
	if signer != address {
		return errInvalidSignature
	}
	return nil
}
