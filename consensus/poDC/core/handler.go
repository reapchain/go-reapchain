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
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/poDC"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
)

// Start implements core.Engine.Start

// proposer = front node in podc
func (c *core) Start(lastSequence *big.Int, lastProposer common.Address, lastProposal poDC.Proposal) error {
	// Initialize last proposer
	c.lastProposer = lastProposer
	c.lastProposal = lastProposal
	c.valSet = c.backend.Validators(c.lastProposal)

	// Start a new round from last sequence + 1
	// a New Round of the State Machine of PoDC : yichoi
	c.startNewRound(&poDC.View{
		Sequence: new(big.Int).Add(lastSequence, common.Big1),
		Round:    common.Big0,
	}, false)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents() /* 이벤트를 수신하는 이벤트 핸들러 */
	go c.handleEvents() /* 이벤트를 처리하는 핸들러로 분리 되어 있음 */

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.unsubscribeEvents()
	return nil
}

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		poDC.RequestEvent{},
		poDC.MessageEvent{},
		poDC.FinalCommittedEvent{},
		// internal events
		backlogEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
}

func (c *core) handleEvents() {
	for event := range c.events.Chan() {
		// A real event arrived, process interesting content
		switch ev := event.Data.(type) {
		case poDC.RequestEvent:
			r := &poDC.Request{
				Proposal: ev.Proposal,
			}
			err := c.handleRequest(r)
			if err == errFutureMessage {
				c.storeRequestMsg(r)
			}
		case poDC.MessageEvent:
			c.handleMsg(ev.Payload)
		case poDC.FinalCommittedEvent:
			c.handleFinalCommitted(ev.Proposal, ev.Proposer)
		case backlogEvent:
			// No need to check signature for internal messages
			c.handleCheckedMsg(ev.msg, ev.src)
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(payload []byte) error {
	logger := c.logger.New("address", c.address)

	// Decode message and check its signature
	msg := new(message)
	if err := msg.FromPayload(payload, c.validateFn); err != nil {
		logger.Error("Failed to decode message from payload", "err", err)
		return err
	}

	// Only accept message if the address is valid
	_, src := c.valSet.GetByAddress(msg.Address)
	if src == nil {
		logger.Error("Invalid address in message", "msg", msg)
		return istanbul.ErrUnauthorizedAddress
	}

	return c.handleCheckedMsg(msg, src)
}

func (c *core) handleCheckedMsg(msg *message, src istanbul.Validator) error {
	logger := c.logger.New("address", c.address, "from", src)

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		if err == errFutureMessage {
			c.storeBacklog(msg, src)
			return nil
		}

		return err
	}

	switch msg.Code {
	case msgPreprepare:
		return testBacklog(c.handlePreprepare(msg, src))
	case msgPrepare:
		return testBacklog(c.handlePrepare(msg, src))

	case msgD_select: //yichoi
		return testBacklog(c.handlePrepare(msg, src))
	case msgCommit:
		return testBacklog(c.handleCommit(msg, src))

	case msgD_commit: //yichoi d-commit
		return testBacklog(c.handleD_commit(msg, src))

	case msgGetCandiateList:
		return testBacklog(c.handleCandidateList(msg, src))

	case msgStartRacing:
		return testBacklog(c.handleStartRacing(msg, src))

	case msgRegisterCommittee:
		return testBacklog(c.RegisterCommittee(msg, src))

	case msgRoundChange:
		return c.handleRoundChange(msg, src)
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}
