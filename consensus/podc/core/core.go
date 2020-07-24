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
	"bytes"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	goMetrics "github.com/rcrowley/go-metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

// New creates an PoDC consensus core
func New(backend podc.Backend, config *podc.Config) Engine {
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              StateAcceptRequest,
		logger:             log.New("address", backend.Address()),
		backend:            backend,
		backlogs:           make(map[podc.Validator]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},

		roundMeter:         metrics.NewMeter("consensus/podc/core/round"),
		sequenceMeter:      metrics.NewMeter("consensus/podc/core/sequence"),
		consensusTimer:     metrics.NewTimer("consensus/podc/core/consensus"),
		racingMu:			new(sync.Mutex),

	}
	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *podc.Config
	address common.Address
	state   State
	logger  log.Logger

	backend podc.Backend
	events  *event.TypeMuxSubscription


	qmanager common.Address

	lastProposer          common.Address

	lastProposal          podc.Proposal
	lastSequence          *big.Int  //yichoi added for solving inconsistent
	valSet                podc.ValidatorSet

	waitingForRoundChange bool

	validateFn            func([]byte, []byte) (common.Address, error)

	backlogs   map[podc.Validator]*prque.Prque
	backlogsMu *sync.Mutex

	current *roundState

	roundChangeSet   *roundChangeSet
	roundChangeTimer *time.Timer

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	consensusTimestamp time.Time
	// the meter to record the round change rate
	roundMeter goMetrics.Meter
	// the meter to record the sequence update rate
	sequenceMeter goMetrics.Meter
	// the timer to record consensus duration (from accepting a preprepare to final committed stage)
	consensusTimer goMetrics.Timer

	tag common.Tag
	count uint
	racingFlag bool
	racingMu *sync.Mutex
	startTime time.Time
	intervalTime time.Time
	criteria float64   //Criteria 	int		`json:"criteria"` type.go에서는  int 로 받았다가, 여기서는 float64
	ExtraDataLength int
}
// 최종 전송할 메시지를 만듦
func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.Address()

	// Add proof of consensus
	msg.CommittedSeal = []byte{}
	// Assign the CommittedSeal if it's a commit message and proposal is not nil
	if msg.Code == msgCommit && c.current.Proposal() != nil {
		seal := PrepareCommittedSeal(c.current.Proposal().Hash())
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) send(msg *message, target common.Address) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Send payload
	if err = c.backend.Send(payload, target); err != nil {
		logger.Error("Failed to send message", "msg", msg, "err", err)
		return
	}
}
// message 구조체 내에 enode address가 있음.
func (c *core) broadcast(msg *message) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(c.valSet, payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) multicast(msg *message, targets []common.Address) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	if err = c.backend.Multicast(payload, targets); err != nil {
		logger.Error("Failed to multicast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) currentView() *podc.View {
	return &podc.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
	}
}
func (c *core) isRequestQman() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
	         //Front node 인가
}
func (c *core) isProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
}

func (c *core) commit() {
	c.setState(StateCommitted)

	proposal := c.current.Proposal()
	if proposal != nil {
		var signatures []byte
		for _, v := range c.current.Commits.Values() {
			signatures = append(signatures, v.CommittedSeal...)
		}

		if err := c.backend.Commit(proposal, signatures); err != nil {
			c.sendNextRoundChange()
			return
		}
	}
}

func (c *core) startNewRound(newView *podc.View, roundChange bool) {

	var logger log.Logger
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_seq", 0, "old_proposer", c.valSet.GetProposer())
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())  //1.
	}

	c.valSet = c.backend.Validators(c.lastProposal)
	// Clear invalid round change messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.current = newRoundState(newView, c.valSet)
	// Calculate new proposer
	c.valSet.CalcProposer(c.lastProposer, newView.Round.Uint64(), c.qmanager)
	c.waitingForRoundChange = false
	c.setState(StateRequestQman)
	if roundChange && c.isProposer() {
		c.backend.NextRound()
	}
	c.newRoundChangeTimer()
	//logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size())
	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "size", c.valSet.Size())
}

func (c *core) catchUpRound(view *podc.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())

	if view.Round.Cmp(c.current.Round()) > 0 {
		c.roundMeter.Mark(new(big.Int).Sub(view.Round, c.current.Round()).Int64())
	}
	c.waitingForRoundChange = true
	c.current = newRoundState(view, c.valSet)
	c.roundChangeSet.Clear(view.Round)
	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_seq", view.Sequence, "new_proposer", c.valSet)
}

func (c *core) 	setState(state State) {
	if c.state != state {  //상태가 다르면, 입력 파라미터 상태로 설
		c.state = state
	}
	if state == StateAcceptRequest || state == StateRequestQman {
		c.processPendingRequests() //바로 보내지않고, 지연시켜서 보내는, sendevent 핸들러에게 보내서,,,
	}
	c.processBacklog()
}

func (c *core) Address() common.Address {
	return c.address
}

func (c *core) newRoundChangeTimer() {
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}

	// set timeout based on the round number
	timeout := time.Duration(c.config.RequestTimeout)*time.Millisecond + time.Duration(c.current.Round().Uint64()*c.config.BlockPeriod)*time.Second
	c.roundChangeTimer = time.AfterFunc(timeout, func() {
		// If we're not waiting for round change yet, we can try to catch up
		// the max round with F+1 round change message. We only need to catch up
		// if the max round is larger than current round.
		if !c.waitingForRoundChange {  // bool 값
			maxRound := c.roundChangeSet.MaxRound(c.valSet.F() + 1)
			if maxRound != nil && maxRound.Cmp(c.current.Round()) > 0 {
				c.sendRoundChange(maxRound)
			} else {
				c.sendNextRoundChange()
			}
		} else {
			c.sendNextRoundChange()
		}
	})
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return podc.CheckValidatorSignature(c.valSet, data, sig)
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}

func (c *core) Tag() common.Tag {
	return c.tag
}

func (c *core) SetTag(t common.Tag) {
	c.tag = t
}

func (c *core) GetValidatorListExceptQman() []common.Address {
	var addrList []common.Address

	for _, val := range c.valSet.List() {
		if val.Address() != c.qmanager {
			addrList = append(addrList, val.Address())
		}
	}

	return addrList
}
