/* PoDC D-Select */

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
	"reflect"
)

/* D-selct has 6 states.
1. cordinator and candidate validator select

2. c-comfirm to QM

3. C-receive ( candidate list confirm )

4. C-request ( candidate fix and notification )

5. racing

6. d-selected completed



*/

// just I renamed func name in order to implement podc, so yours modify inside of func below

func (c *core) sendD_select() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgPrepare,
		Msg:  encodedSubject,
	})
}

// TODO: define FrontNode like
// type FrontNode struct {
//		address common.Address
// }
func (c *core) handleD_select(msg *message, src poDC.FrontNode) error {
	// Decode prepare message
	var dselect *poDC.Subject
	err := msg.Decode(&dselect)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(msgPrepare, dselect.View); err != nil {
		return err
	}

	if err := c.verifyD_select(dselect, src); err != nil {
		return err
	}

	c.identifySelfNode(msg, src)

	c.acceptD_select(msg, src)

	// Change to StatePrepared if we've received enough prepare messages
	// and we are in earlier state before StatePrepared
	if c.current.Prepares.Size() > 2*c.valSet.F() && c.state.Cmp(StatePrepared) < 0 && c.current.D_select.selfNode.isCoordinator() {
		// c.setState(StatePrepared)
		c.requestCandidateList()
	}

	return nil
}

// verifyPrepare verifies if the received prepare message is equivalent to our subject
func (c *core) verifyD_select(prepare *poDC.Subject, src poDC.FrontNode) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logger.Warn("Inconsistent subjects between prepare and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptD_select(msg *message, src poDC.FrontNode) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the prepare message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logger.Error("Failed to add prepare message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}

// identify what self-node kind is : committee candidate, coordinator, senator
func (c *core) identifySelfNode(msg *message, src poDC.FrontNode) error {
	// TODO
}

// send request to QManager to get committee candidate list
func (c *core) requestCandidateList() error {
	// TODO
}

// when QManager response, handle the candidate list and notify selection of coordinator to committee candidates
func (c *core) handleCandidateList(msg *message, src poDC.QManager) error {
	// TODO

	// target would be c.current.D_select.Coordinator.candidates
	c.notifyCoordinator()
}

// if a committee candidate listen notification from coordinator, start racing
func (c *core) handleStartRacing(msg *message, src poDC.Coordinator) error {
	// TODO

	// target would be c.current.D_select.Coordinator
	c.racing()
}

// register committee candidate to committee by the order of arrival
// and set round as d-commit
func (c *core) handleRegisterCommittee(msg *messgae, src poDC.Candidate) error {
	// TODO

	if committee.Number() == 15 {
		c.setState(StateD_selected)
		c.sendD_commit()
	}
}
