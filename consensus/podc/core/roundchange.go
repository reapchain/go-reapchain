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
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc"
)

// sendNextRoundChange sends the round change message with current round + 1
func (c *core) sendNextRoundChange() {

	cv := c.currentView()
	c.sendRoundChange(new(big.Int).Add(cv.Round, common.Big1))

}

// sendRoundChange sends the round change message with the given round
func (c *core) sendRoundChange(round *big.Int) {
	logger := c.logger.New("state", c.state)

	cv := c.currentView()
	if cv.Round.Cmp(round) >= 0 {
		logger.Error("Cannot send out the round change", "current round", cv.Round, "target round", round)
		return
	}

	c.catchUpRound(&podc.View{
		// The round number we'd like to transfer to.
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(cv.Sequence),
	})

	// Now we have the new round number and sequence number
	cv = c.currentView()
	rc := &roundChange{
		Round:    new(big.Int).Set(cv.Round),
		Sequence: new(big.Int).Set(cv.Sequence),
		Digest:   common.Hash{},
	}

	payload, err := Encode(rc)
	if err != nil {
		logger.Error("Failed to encode round change", "rc", rc, "err", err)
		return
	}

	c.broadcast(&message{
		Code: msgRoundChange,
		Msg:  payload,
	})
}

func (c *core) handleRoundChange(msg *message, src podc.Validator) error {
	logger := c.logger.New("state", c.state)
		// Decode round change message
		var rc *roundChange
	if( !podc_global.QManConnected ) { // if i'm not Qman and general geth. then roundchange and start new round. for qman, don't roundchange, it is not necessary.

		if err := msg.Decode(&rc); err != nil {
			logger.Error("Failed to decode round change", "err", err)
			return errInvalidMessage
		}

		cv := c.currentView()
		//log.Info("I'm  not the Qmanager(handleRoundChange) ")

		// We never accept round change message with different sequence number
		if rc.Sequence.Cmp(cv.Sequence) != 0 {
			logger.Warn("Inconsistent sequence number(handleRoundChange)", "expected", cv.Sequence, "got", rc.Sequence)
            // sequence no가 틀리면, 블럭 싱크 다시 맞춰주기 위해서, 싱크로나이즈 동작 수행 해주면 됨. 여기서,
            // 싱크로나이즈를 수행하기 위해서 여기서 싱크코드 넣고, 시퀀스 맞춰지면,, 다시 startNewRound. 수행
            // 메시지 fail 나는 경우, remove peer를 하기 때문에 쓸모없는 peer로,, 에러메시지 나오고, EOF,, 등등.
            // 싱크 에러와, fetcher를,, 여기서 채널로,, fetcher or sync hanlder로,, 채널을 통해서 메시지르 보내면,
            // eth/handler.go 의 synloop()에서,, 처리되게할 것.

			//c.startNewRound(&podc.View{
			//	Round:    new(big.Int).Set(rc.Round),
			//	Sequence: new(big.Int).Set(rc.Sequence),
			//}, true)

			return errInvalidMessage
		}

		// We never accept round change message with smaller round number
		if rc.Round.Cmp(cv.Round) < 0 {
			logger.Warn("Old round change", "from", src, "expected", cv.Round, "got", rc.Round)
			return errOldMessage
		}

		num, err := c.roundChangeSet.Add(rc.Round, msg)

		if err != nil {
			logger.Warn("Failed to add round change message", "from", src, "msg", msg, "err", err)
			return err
		}

		// Once we received f+1 round change messages, those messages form a weak certificate.
		// If our round number is smaller than the certificate's round number, we would
		// try to catch up the round number.
		if c.waitingForRoundChange && num == int(c.valSet.F()+1) {
			if cv.Round.Cmp(rc.Round) < 0 {
				c.sendRoundChange(rc.Round)
			}
		}

		// We've received 2f+1 round change messages, start a new round immediately.
		if num == int(2*c.valSet.F()+1) {
			c.startNewRound(&podc.View{
				Round:    new(big.Int).Set(rc.Round),
				Sequence: new(big.Int).Set(rc.Sequence),
			}, true)
		}

		//}
	}

	return nil
}

// ----------------------------------------------------------------------------

func newRoundChangeSet(valSet podc.ValidatorSet) *roundChangeSet {
	return &roundChangeSet{
		validatorSet: valSet,
		roundChanges: make(map[uint64]*messageSet),
		mu:           new(sync.Mutex),
	}
}

type roundChangeSet struct {
	validatorSet podc.ValidatorSet
	roundChanges map[uint64]*messageSet
	mu           *sync.Mutex
}

// Add adds the round and message into round change set
func (rcs *roundChangeSet) Add(r *big.Int, msg *message) (int, error) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	round := r.Uint64()
	if rcs.roundChanges[round] == nil {
		rcs.roundChanges[round] = newMessageSet(rcs.validatorSet)
	}
	err := rcs.roundChanges[round].Add(msg)  //?
	if err != nil {
		return 0, err
	}
	return rcs.roundChanges[round].Size(), nil
}

// Clear deletes the messages with smaller round
func (rcs *roundChangeSet) Clear(round *big.Int) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	for k, rms := range rcs.roundChanges {
		if len(rms.Values()) == 0 || k < round.Uint64() {
			delete(rcs.roundChanges, k)
		}
	}
}

// MaxRound returns the max round which the number of messages is equal or larger than num
func (rcs *roundChangeSet) MaxRound(num int) *big.Int {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	var maxRound *big.Int
	for k, rms := range rcs.roundChanges {
		if rms.Size() < num {
			continue
		}
		r := big.NewInt(int64(k))
		if maxRound == nil || maxRound.Cmp(r) < 0 {
			maxRound = r
		}
	}
	return maxRound
}
