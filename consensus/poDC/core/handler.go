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
	// "github.com/ethereum/go-ethereum/consensus/istanbul"
)

// Start implements core.Engine.Start

// proposer = front node in podc
// 앞에서 Proposer가 결정된 상태에서, 여기서 합의 시작
func (c *core) Start(lastSequence *big.Int, lastProposer common.Address, lastProposal poDC.Proposal) error {
	// Initialize last proposer
	c.lastProposer = lastProposer  // 여기서 lastProposer는 최신을 의미
	c.lastProposal = lastProposal  //합의할 블럭을 제시하는데, 시리얼화된 블럭을 데이터로 가져옴.
	c.valSet = c.backend.Validators(c.lastProposal)

	// Start a new round from last sequence + 1
	// a New Round of the State Machine of PoDC : yichoi
	c.startNewRound(&poDC.View{ // 1. 최초 라운드 시작
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
		poDC.RequestEvent{}, // Qmanger request event도 여기.. 속함.
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
// 새 라운드에서, Proposer는 블럭을 노드에 던져주면 끝나고, D-select, D-commit 상태는,, 랜덤으로 선정된. 코디와 큐매니저끼리 검증하고,
// 최종 검증된 블럭을 다시 전체 노드에 던지면, 그때,, 이 노드가,, 타이머 등을 리셋하고,
// 검증된 블럭을 자신의 노드의 체인에 연결해서 새 블럭을 생성시키면 된다. ?

// 여기서 블럭 전체를 노드에 전파해야,, 코디가,, 해당블럭을 검증하고, 커밋을 완료하게 된다.

func (c *core) handleEvents() {
	for event := range c.events.Chan() {
		// A real event arrived, process interesting content
		switch ev := event.Data.(type) {
		case poDC.RequestEvent:  // Proposal 구조체
			r := &poDC.Request{
				Proposal: ev.Proposal,
			}
			err := c.handleRequest(r)
			if err == errFutureMessage {
				c.storeRequestMsg(r)
			}
		case poDC.MessageEvent:  //cordi로부터 합의를 통해서, 최종 검증된 블럭을 받으면, payload  배열
			c.handleMsg(ev.Payload)
		case poDC.FinalCommittedEvent:  //Proposal 구조체
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
		return poDC.ErrUnauthorizedAddress
	}
	// 위는 코어 쪽에서 메시지 받고, 검증하고, Payload를 가져오고,
	return c.handleCheckedMsg(msg, src) // 메시지 처리 핸들러는 여기서,
}

// 실제 모든 코어단에서 받는 메시지를 처리하는 마지막 핸들러...

func (c *core) handleCheckedMsg(msg *message, src poDC.Validator) error {
	logger := c.logger.New("address", c.address, "from", src) // 송진자 enode address를 가지고, 로깅

	// Store the message if it's a future message
	testBacklog := func(err error) error {
		if err == errFutureMessage {
			c.storeBacklog(msg, src)
			return nil
		}

		return err
	}

	switch msg.Code {

	// 경우에 따라서 테스트시 오류 검증을 위해서, 이스탄불 소스를 돌릴수 있기 때문에, 여기서는 메시지 상태 머신을 섞어서 쓴다.
	// Qmanager로부터 메시지를 받으면 블럭헤더의 ExtraDATA의 내용을 채운다.
	// 상원, 하원, 운영위 후보군, 코디 등.
	case msgReceivedFromQman:
		return testBacklog(c.handleExtraData(msg, src))  //handleExtraData 여기서 Extradata를 블럭헤더에 넣어서처리

	case msgPreprepare:
		return testBacklog(c.handlePreprepare(msg, src))

	/*
	case msgPrepare:
		return testBacklog(c.handlePrepare(msg, src))

	case msgCommit:
		return testBacklog(c.handleCommit(msg, src))
    */
		// message handler - begin
	case msgD_select: //yichoi
		return testBacklog(c.handlePrepare(msg, src))
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
// cordi로부터  최종 확정된 블럭을 받으면, 커밋
	case msgGetVerifiedBlock:
		return testBacklog(c.handleCommit(msg, src))
//

	default:
		logger.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}
