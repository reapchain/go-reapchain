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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
)

//yichoi - begin
// QmanagerNodes returns a list of node enode URLs configured as Qmanager nodes.
//func (c *Config) QmanagerNodes() []*discover.Node {
//	return c.parsePersistentNodes(c.resolvePath(datadirQmanagerNodes))
//}
// end
/*func( c *core) getQman_enode() []*discover.Node {

	var d node.p2p.c.qmanager
	d.p2p.Qmanagernodes = node.d.QmanagerNodes( )

} */
/* 최초 Qmanager 에게 ExtraDATA를 요청하는 단계 */
/*
func (c *core) sendRequestExtraDataToQman(request *istanbul.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() { //?
		curView := c.currentView()
		preprepare, err := Encode(&istanbul.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		// load Qmanager enode address from static node

		config := &p2p.Config{ Name: "unit-test", DataDir: dir, P2P: p2p.Config{PrivateKey: key}
			config.QmanagerNodes.


			// make request massage,
			/* type message struct {
				Code          uint64
				Msg           []byte
				Address       common.Address
				Signature     []byte
				CommittedSeal []byte
			} */
	/*		var payload = makemsg()
			var d Config
			var qman_enode common.Address
			qman_enode  = getQman_enode( )
			qman_enode =""

			c.send(&message{
			Code: msgRequestQman,
			Msg: payload
			Address:  qman_enode ,
		})

			// proposal block 전파는 핸들러로 옮겨야,, Qmanager에서 수신시,, 처리되게끔.  / pre-prepare 상태
			// 다음은 d-select 상태로 상태 전이함.
		}
	}
	*/
func (c *core) sendPreprepare(request *istanbul.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.isProposer() {
		curView := c.currentView()
		preprepare, err := Encode(&istanbul.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}
		var account_addr common.Address  //ethereum account : 20 byte
		//var payload []byte
		//payload = ""
		//qman_enode = "e81bd88b5c3a9a7eebb454eb3fab0988d2134ef2fa3066b5b40f8719a44cce52c032b1dc698e28571cf4df95fc6ea386d6eb10ab6c26e91e315334e11175563e@192.168.0.100:30301"
		copy(account_addr[:],"2259172c57bde543b819d21805ad6c3d06e89d18")

		/*type message struct {
			Code          uint64
			Msg           []byte
			Address       common.Address
			Signature     []byte
			CommittedSeal []byte
		} */
		c.send(&message{
			Code: msgRequestQman,
			Msg: preprepare,
			Address:  account_addr,
			//Signature: 1,
			//CommittedSeal: 1,
		}, account_addr )

		c.broadcast(&message{
			Code: msgPreprepare,
			Msg:  preprepare,
		})


	}
}

func (c *core) handlePreprepare(msg *message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode preprepare
	var preprepare *istanbul.Preprepare
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

	// Verify the proposal we received
	if err := c.backend.Verify(preprepare.Proposal); err != nil {
		logger.Warn("Failed to verify proposal", "err", err)
		c.sendNextRoundChange()
		return err
	}

	if c.state == StateAcceptRequest {
		c.acceptPreprepare(preprepare)
		c.setState(StatePreprepared)
		c.sendPrepare()
	}

	return nil
}

func (c *core) acceptPreprepare(preprepare *istanbul.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}
