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
	"github.com/ethereum/go-ethereum/consensus/poDC"
	"time"

//	"github.com/ethereum/go-ethereum/consensus/istanbul"
)

/* 최초 Qmanager 에게 ExtraDATA를 요청하는 단계 */
func (c *core) sendRequestExtraDataToQman(request *poDC.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isRequestQman() { //?
		curView := c.currentView()
		preprepare, err := Encode(&poDC.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		// proposal block 전파 / pre-prepare 상태

		/* c.broadcast(&message{
			Code: msgPreprepare,
			Msg:  preprepare,
		}) */
		// 다음은 d-select 상태로 상태 전이함.
	}
}
func (c *core) sendPreprepare(request *poDC.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() {
		curView := c.currentView()
		preprepare, err := Encode(&poDC.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

// proposal block 전파 / pre-prepare 상태

		c.broadcast(&message{
			Code: msgPreprepare,
			Msg:  preprepare,
		})
		// 다음은 d-select 상태로 상태 전이함.
	}
}

func (c *core) handlePreprepare(msg *message, src poDC.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode preprepare
	var preprepare *poDC.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return errFailedDecodePreprepare
	}

	// Ensure we have the same view with the preprepare message
	if err := c.checkMessage(msgPreprepare, preprepare.View); err != nil {
		return err
	}

	// Check if the message comes from current proposer( = Front node ) in PoDC
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore preprepare messages from non-proposer")
		return errNotFromProposer
	}

	// Verify the proposal we received
	if err := c.backend.Verify(preprepare.Proposal); err != nil {
		logger.Warn("Failed to verify proposal", "err", err)
		c.sendNextRoundChange()
		return err
	}

	// 상태 전이 모드
	if c.state == StateAcceptRequest {
		c.acceptPreprepare(preprepare)
		c.setState(StatePreprepared)
		c.sendPrepare()
	}

	return nil
}

func (c *core) acceptPreprepare(preprepare *poDC.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}
/* begin yichoi added */
func (c *core) acceptD_select(preprepare *poDC.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetD_select(d_select)
}
func (c *core) acceptD_commit(preprepare *poDC.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetD_commit(d_commit)
}
/* end */
