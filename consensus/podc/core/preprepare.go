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
	"github.com/ethereum/go-ethereum/common"
	"time"

	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) sendRequestExtraDataToQman(request *podc.Request) {
	logger := c.logger.New("state", c.state)
	log.Info("2. Interval time between start new round and pre-prepare", "elapsed", common.PrettyDuration(time.Since(c.startTime)))
	c.intervalTime = time.Now()

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() { //?
		curView := c.currentView()
		preprepare, err := Encode(&podc.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}


		if c.valSet.IsProposer(c.Address()) {  //Proposer인 경우만 Qman에게 요청 메시지 보낸다.
			c.broadcast(&message{
				Code: msgRequest,  //request to qman to get ExtraData
				Msg: preprepare,
				Address: c.qmanager,
			})
		}
			// proposal block 전파는 핸들러로 옮겨야,, Qmanager에서 수신시,, 처리되게끔.  / pre-prepare 상태
			// 다음은 d-select 상태로 상태 전이함.

	}
}
// 2. go to step 2 : pre-prepare step
func (c *core) sendPreprepare(request *podc.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() {
		curView := c.currentView()
		preprepare, err := Encode(&podc.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&message{
			Code: msgPreprepare,
			Msg:  preprepare,
		})


	}
}


func (c *core) handleQmanager(msg *message, src podc.Validator) error {  //request to qman
	logger := c.logger.New("from", src, "state", c.state)
	//if (!reflect.DeepEqual(c.qmanager, c.Address())) {
		log.Info("I'm not Qman in handleQmanager" )

		// Qmanager receiver에 맞게 수정할 부분 begin
		// 1. Extra data 전송하고,
		// 2. Enrollment 하고, martin
		// Cordi가 "자신기 코디"임을 보내오면,
		// Cordi에게 C-Confirm 를 보내고,

		// Decode preprepare
		var preprepare *podc.Preprepare
		err := msg.Decode(&preprepare)
		if err != nil {
			return errFailedDecodePreprepare
		}

		// Ensure we have the same view with the preprepare message
		if err := c.checkMessage(msgPreprepare, preprepare.View); err != nil {
			return err
		}

		// Check if the message comes from current proposer
		if !c.valSet.IsProposer(src.Address()) {
			logger.Warn("Ignore preprepare messages from non-proposer")
			return errNotFromProposer
		}

		if c.valSet.IsProposer(c.Address()) {  // I'm Front node.
			log.Info("I'm Proposer!!!!!!!")
		}
		// Verify the proposal we received
		if err := c.backend.Verify(preprepare.Proposal); err != nil {
			logger.Warn("handleQmanager: Failed to verify proposal", "err", err) //?
			c.sendNextRoundChange()                                              //important : inconsistent mismatch ...
			return err
		}

		log.Info("3. Set pre-prepare state", "elapsed", common.PrettyDuration(time.Since(c.intervalTime)))
		c.intervalTime = time.Now()

		if c.state == StateRequest {
			c.acceptPreprepare(preprepare)
			c.setState(StatePreprepared)
			//c.sendPrepare()
			if c.valSet.IsProposer(c.Address()) {
				c.sendExtraDataRequest()
			}
		}
		// 수정할 부분 end
	//}
	return nil
}

func (c *core) handlePreprepare(msg *message, src podc.Validator) error{
		logger := c.logger.New("from", src, "state", c.state)

		// Decode preprepare
		var preprepare *podc.Preprepare
		err := msg.Decode(&preprepare)
		if err != nil{
		return errFailedDecodePreprepare
	}

		// Ensure we have the same view with the preprepare message
		if err := c.checkMessage(msgPreprepare, preprepare.View); err != nil{
		return err
	}

		// Check if the message comes from current proposer
		if !c.valSet.IsProposer(src.Address()){
		logger.Warn("Ignore preprepare messages from non-proposer")
		return errNotFromProposer
	}

		if c.valSet.IsProposer(c.Address()){
			log.Info("I'm Proposer!!!!!!!")
	}
		// Verify the proposal we received
		if err := c.backend.Verify(preprepare.Proposal); err != nil{
		logger.Warn("Failed to verify proposal", "err", err)
		c.sendNextRoundChange()
		return err
	}

		if c.state == StateAcceptRequest{
		c.acceptPreprepare(preprepare)
		c.setState(StatePreprepared)
		log.Info("2. handlePreprepare: StatePreprepared")
		if c.valSet.IsProposer(c.Address()){
		c.sendDSelect()
	}

	}
	return nil
}

func (c *core) acceptPreprepare(preprepare *podc.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}