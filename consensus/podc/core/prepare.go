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
	"reflect"

	"github.com/ethereum/go-ethereum/consensus/podc"
)

func (c *core) sendPrepare() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{  // Proposal Block 전파
		Code: msgPrepare,
		Msg:  encodedSubject,
	})
}
//
//func (c *core) handlePrepare(msg *message, src podc.Validator) error {
//	// Decode prepare message
//	var prepare *podc.Subject
//	err := msg.Decode(&prepare)
//	if err != nil {
//		return errFailedDecodePrepare
//	}
//
//	if err := c.checkMessage(msgPrepare, prepare.View); err != nil {
//		return err
//	}
//
//	if err := c.verifyPrepare(prepare, src); err != nil {
//		return err
//	}
//
//	c.acceptPrepare(msg, src)
//
//	// Change to StatePrepared if we've received enough prepare messages
//	// and we are in earlier state before StatePrepared
//	if c.current.Prepares.Size() > 2*c.valSet.F() && c.state.Cmp(StatePrepared) < 0 {
//		c.setState(StatePrepared)
//		c.sendCommit()
//	}
//
//	return nil
//}

// verifyPrepare verifies if the received prepare message is equivalent to our subject
func (c *core) verifyPrepare(prepare *podc.Subject, src podc.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logger.Warn("Inconsistent subjects between prepare and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptPrepare(msg *message, src podc.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the prepare message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logger.Error("Failed to add prepare message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
