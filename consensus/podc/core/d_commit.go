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
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
	"reflect"
	"time"
)

func (c *core) sendDCommit() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	if sub == nil {
		logger.Error("Failed to get Subject")
		return
	}
	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgDCommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handleDCommit(msg *message, src podc.Validator) error {
	// Decode commit message
	var commit *podc.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(msgDCommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyDCommit(commit, src); err != nil {
		return err
	}

	c.acceptDCommit(msg, src)

	// Commit the proposal once we have enough commit messages and we are not in StateCommitted.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without prepare messages.
	if c.current.Dcommits.Size() > 2*c.valSet.F() && c.state.Cmp(StateDCommitted) < 0 {
		c.commit()
		log.Info("6. D-commit end", "elapsed", common.PrettyDuration(time.Since(c.intervalTime)))
		log.Info("Total Time", "elapsed", common.PrettyDuration(time.Since(c.startTime)))
	}

	return nil
}

// verifyCommit verifies if the received commit message is equivalent to our subject
func (c *core) verifyDCommit(commit *podc.Subject, src podc.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal(verifyDCommit)", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptDCommit(msg *message, src podc.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the commit message to current round state
	if err := c.current.Dcommits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	return nil
}
