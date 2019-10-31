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
//	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/consensus/poDC"  //yichoi
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	metrics "github.com/ethereum/go-ethereum/metrics"
	goMetrics "github.com/rcrowley/go-metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

// New creates an PoDC consensus core
func New(backend poDC.Backend, config *poDC.Config) Engine {
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              StateAcceptRequest,
		logger:             log.New("address", backend.Address()),
		backend:            backend,
		backlogs:           make(map[poDC.Validator]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},
		roundMeter:         metrics.NewMeter("consensus/poDC/core/round"),
		sequenceMeter:      metrics.NewMeter("consensus/poDC/core/sequence"),
		consensusTimer:     metrics.NewTimer("consensus/poDC/core/consensus"),
	}
	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *poDC.Config
	address common.Address
	state   State
	logger  log.Logger

	backend poDC.Backend
	events  *event.TypeMuxSubscription

	lastProposer          common.Address
	lastProposal          poDC.Proposal
	valSet                poDC.ValidatorSet
	waitingForRoundChange bool
	waitingForStateChange bool // added for PoDC, because Proposer release it's right to random coordinator
    //statechange flag will be need, I think , by yichoi

	validateFn            func([]byte, []byte) (common.Address, error)

	backlogs   map[poDC.Validator]*prque.Prque
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
}
// 최종 전송할 메시지를 만듦
func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.Address()  //message 에 송신자 enode 주소를 탑재

	// Add proof of consensus
	msg.CommittedSeal = []byte{}  // CommittedSeal 배열 초기화
	// Assign the CommittedSeal if it's a commit message and proposal is not nil
	if msg.Code == msgCommit && c.current.Proposal() != nil {
		seal := PrepareCommittedSeal(c.current.Proposal().Hash())  // message 구조체에 CommittedSeal 배열을 채움
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

	payload, err := c.finalizeMessage(msg)  //최종적으로 메시지 구조체에 탑재할 모든 메시지를 만듦
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

func (c *core) currentView() *poDC.View {
	return &poDC.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
	}
}

/* Qmanager에게 ExtraData를 요청하는 거 맞는지 확인 하는 함수 */
/*
func (c *core) isRequestQman() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsRequestQman(c.backend.Address())
}
*/


func (c *core) isProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())  //Proposer인지 체크함. 여기서 ,
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
// state machine 의 NewRound start 여기서  .. yichoi
func (c *core) startNewRound(newView *poDC.View, roundChange bool) {
	var logger log.Logger

	// Qmanager에게 Extradata 요청

    // 핸들러는 나중에 추가하고,
    // 여기서는 일단, Qman 에서 바로 응답이 온다고 가정하고, 짜고, 나중에 핸들러로 옮긴다. 2단계로.

   //1.  send qman

   //2. receive qman

   //3. 그냥 블럭을 모든 노드에 던지면, 블럭연결 제안자는 여기서 빠지고,
       // 임의로 선정된 코디가,, Qman에 통보하고, 자기가, 그 받은 블럭을 검증하고, 처리해서,
       // 체인에 연결하고, 전체 노드에 전파시,
       // 내 노드가 핸들러를 통해서, 받는다.
       // 핸들러에서 처리되면,
       // 아래 타이머를 리셋해서,, 본 라운드가 끝났다는 것을 통지하고, 다음 라운드로 돌아간다.
       // Proposer( = Front node in PoDC ) 가,, 자기 권한을 임의로 선정된 코디에게 위임하는 개념잉
















	/* if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_seq", 0, "old_proposer", c.valSet.GetProposer())
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())
	}



	c.valSet = c.backend.Validators(c.lastProposal) */



	// Clear invalid round change messages
/* 	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.current = newRoundState(newView, c.valSet)
	// Calculate new proposer
	c.valSet.CalcProposer(c.lastProposer, newView.Round.Uint64())

 */
	c.waitingForRoundChange = false
	c.waitingForStateChange = true   //


	c.setState(StateAcceptRequest)
	/* if roundChange && c.isProposer() {
		c.backend.NextRound()
	} */
	if roundChange  {
		c.backend.NextRound()
	}
	c.newRoundChangeTimer()  //마냥 기다릴수 없어서, 여기서 타이머 설정하고, 빠져나감, 나머지는 메시지/이벤트 핸들러에서, 이벤트/메시지 수신시
	// 처리함.

	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size())
}

func (c *core) catchUpRound(view *poDC.View) {
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

func (c *core) setState(state State) {
	if c.state != state {
		c.state = state
	}
	if state == StateAcceptRequest {
		c.processPendingRequests()
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
	           // 타임아웃 시간은 우측 수식에 의해서 계산됨 값.
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
	return poDC.CheckValidatorSignature(c.valSet, data, sig)
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}
