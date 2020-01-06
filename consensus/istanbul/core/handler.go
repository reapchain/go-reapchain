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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"

	//"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// Start implements core.Engine.Start
func (c *core) Start(lastSequence *big.Int, lastProposer common.Address, lastProposal istanbul.Proposal, qmanager []*discover.Node) error {
	// Initialize last proposer
	c.lastProposer = lastProposer
	var err error
	if( qmanager == nil ) {
		return err
	}
    if (len(qmanager) <= 0)  {
    	log.Debug("Qmanager node is not exist")
        return nil
	}
	//	a[low:high]
	/* a := [5]string{"C", "C++", "Java", "Python", "Go"}

	slice1 := a[1:4]
	slice2 := a[:3]
	slice3 := a[2:]
	slice4 := a[:] */
	// nodekey  : 0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036
	// boot.key : 0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036
	//Qmanager enode: 5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c
    QmanEnode := qmanager[0].ID[:]  //여기까지 정상

	c.qmanager = crypto.PublicKeyBytesToAddress(QmanEnode) //common.Address output from this [account addr]              //slice ->
	                                                       //Qmanager account address(20byte): 926ea01d982c8aeafab7f440084f90fe078cba92
	c.lastProposal = lastProposal
	c.valSet = c.backend.Validators(c.lastProposal)  // Validator array 관리


	// Start a new round from last sequence + 1
	start :=time.Now()
	log.Info("start time of consensus of core engine start()", start )
	c.startNewRound(&istanbul.View{
		Sequence: new(big.Int).Add(lastSequence, common.Big1),
		Round:    common.Big0,
	}, false)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	go c.handleEvents()

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
		istanbul.RequestEvent{},
		istanbul.MessageEvent{},
		istanbul.FinalCommittedEvent{},
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
		case istanbul.RequestEvent:
			r := &istanbul.Request{
				Proposal: ev.Proposal,
			}
			err := c.handleRequest(r)  //send qman here
			if err == errFutureMessage {
				c.storeRequestMsg(r)
			}
		case istanbul.MessageEvent:
			c.handleMsg(ev.Payload)
		case istanbul.FinalCommittedEvent:
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
	/* Qmanager handler for receiving from geth : sending qmanager event */
	case msgHandleQman:
		return testBacklog(c.handleQmanager(msg, src))  //Qmanager receiving event hadler
	case msgPreprepare:
		return testBacklog(c.handlePreprepare(msg, src))
	case msgDSelect:
		return testBacklog(c.handleDSelect(msg, src))
	case msgCoordinatorDecide:
		return testBacklog(c.handleCoordinatorDecide(msg, src))
	case msgRacing:
		return testBacklog(c.handleRacing(msg, src))
	case msgCandidateDecide:
		return testBacklog(c.handleCandidateDecide(msg, src))
	case msgPrepare:
		return testBacklog(c.handlePrepare(msg, src))
	case msgCommit:
		return testBacklog(c.handleCommit(msg, src))
	case msgRoundChange:
		return c.handleRoundChange(msg, src)
	case msgRequestQman:
		return c.handleExtraData(msg, src)
	case msgExtraDataSend:
		return testBacklog(c.handleSentExtraData(msg, src))
	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}
