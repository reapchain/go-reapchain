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
	"github.com/ethereum/go-ethereum/consensus/podc"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	// msgPriority is defined for calculating processing priority to speedup consensus
	// msgPreprepare > msgCommit > msgDSelect
	msgPriority = map[uint64]int{
		msgPreprepare: 1,
		msgCommit:     2,
		msgDSelect:    3,
		//msgDSelect:    3,

	}
)

// checkMessage checks the message state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *podc.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return errInvalidMessage
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}

	if c.waitingForRoundChange {
		return errFutureMessage
	}

	// StateAcceptRequest only accepts msgPreprepare
	// other messages are future messages
	if c.state == StateAcceptRequest {
		if msgCode > msgPreprepare {
			return errFutureMessage
		}
		return nil
	}

	// For states(StatePreprepared, StatePrepared, StateCommitted),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) storeBacklog(msg *message, src podc.Validator) {
	logger := c.logger.New("from", src, "state", c.state)

	if src.Address() == c.Address() {
		logger.Warn("Backlog from self")
		return
	}

	logger.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	backlog := c.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
	}
	switch msg.Code {
	case msgPreprepare:
		var p *podc.Preprepare
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.View))
		}
		// for istanbul.msgDSelect and istanbul.MsgCommit cases
	default:
		var p *podc.Subject
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.View))
		}
	}
	c.backlogs[src] = backlog
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}

		logger := c.logger.New("from", src, "state", c.state)
		isFuture := false

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !(backlog.Empty() || isFuture) {
			m, prio := backlog.Pop()
			msg := m.(*message)
			var view *podc.View
			switch msg.Code {
			case msgPreprepare:  //?
				//msgDSelect and msgDCommit 구현
			case 	msgDSelect:
			case 	msgCommit:
			//코디와 주고 받는것,
			case 	msgCoordinatorDecide:  //1. 코디 결정
			case 	msgRacing:             //2. 레이싱
			case 	msgCandidateDecide:    //3. 후보군 결정


			//case 	msgRoundChange:
			//Qmanager와 주고 받는 것 :
			case 	msgHandleQman:
			case 	msgExtraDataRequest:
			case 	msgExtraDataSend:
			case 	msgCoordinatorConfirmRequest:
			case 	msgCoordinatorConfirmSend:

				var m *podc.Preprepare
				err := msg.Decode(&m)
				if err == nil {
					view = m.View   //메시지뷰를 뷰에 설
				}
			default:
				var sub *podc.Subject
				err := msg.Decode(&sub)
				if err == nil {
					view = sub.View //주제 뷰를 설정
				}
			}
			if view == nil {
				logger.Debug("Nil view", "msg", msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, view)
			if err != nil {
				if err == errFutureMessage {
					logger.Trace("Stop processing backlog", "msg", msg)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				logger.Trace("Skip the backlog event", "msg", msg, "err", err)
				continue
			}
			logger.Trace("Post backlog event", "msg", msg)

			go c.sendEvent(backlogEvent{
				src: src,
				msg: msg,
			})
		}
	}
}

func toPriority(msgCode uint64, view *podc.View) float32 {  //나중에 시퀀스,, 중요. 속도 계산해서,, 추후 처리할 것.
	// FIXME: round will be reset as 0 while new sequence
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Sequence limits the range of round is from 0 to 99
	return -float32(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
}
