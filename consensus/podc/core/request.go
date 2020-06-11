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
	"github.com/ethereum/go-ethereum/log"
	"reflect"
)


func (c *core) handleRequest(request *podc.Request) error {
	logger := c.logger.New("state", c.state, "c.current.sequences", c.current.sequence)

	if err := c.checkRequestMsg(request); err != nil {
		logger.Warn("unexpected requests", "err", err, "request", request)
		return err
	}

	logger.Trace("handleRequest", "request", request.Proposal.Number())

	if c.state ==StateRequest {
		// c.sendPre_prepare()  //send to Qman to get Extradata
		log.Info("StateRequestQman", "StateRequestQman", StateRequest)
		//Qmanager가 아니라면,
		if (!reflect.DeepEqual(c.qmanager, c.Address())) { //if I'm not Qmanager
			log.Info("I'm not Qman  in handleRequest: sendRequestExtraDataToQman" )
			c.sendRequestExtraDataToQman(request)
		}

	}

		//Qmanager response check is more prefer then StateAccepRequest state.
	if c.state == StateAcceptRequest {
		log.Info("StateAcceptRequest", "StateAcceptRequest", StateAcceptRequest)
		c.sendPreprepare(request)
	}





	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
// 리퀘스트 메시지의 에러를 먼저 체크한다..
func (c *core) checkRequestMsg(request *podc.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *podc.Request) {
	logger := c.logger.New("state", c.state)

	logger.Trace("Store future requests", "request", request)

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, float32(-request.Proposal.Number().Int64()))
}
//PendingRequest 는 지연 요청으로, 곧바로 보내지않고, 고루틴을 써서, 지연 처리.. ?
func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()

		r, ok := m.(*podc.Request)

		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "request", r)
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "request", r, "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "request", r)

        /* **************** */
		go c.sendEvent(podc.RequestEvent{ // 여기서 보냄. Qmanager로 보내야함.

			Proposal: r.Proposal,
		})
	}
}


//PendingRequest 는 지연 요청으로, 곧바로 보내지않고, 고루틴을 써서, 지연 처리.. ?
func (c *core) processPendingRequestsQman() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()  //stack에서 우선순위 가져온다.

		r, ok := m.(*podc.Request )
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "request", r)
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "request", r, "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "request", r)
        //data := makemsg()
        //var sample_msg []byte = {1234}
        enode_slice := c.qmanager[:]
		go c.sendEvent(podc.QmanDataEvent{  // 여기서 보냄. Qmanager로 보내야함.
			Target : c.qmanager, //?
			Data : enode_slice , //?

		})
	}
}
