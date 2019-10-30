/* PoDC D-Commit */

package core

import (
	"github.com/ethereum/go-ethereum/consensus/poDC"
	"reflect"
)

func (c *core) sendD_Commit() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgCommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handleD_Commit(msg *message, src poDC.Validator) error {
	// Decode commit message
	var commit *poDC.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(msgCommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyCommit(commit, src); err != nil {
		return err
	}

	c.acceptCommit(msg, src)

	// Commit the proposal once we have enough commit messages and we are not in StateCommitted.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without prepare messages.
	if c.current.Commits.Size() > 2*c.valSet.F() && c.state.Cmp(StateCommitted) < 0 {
		c.commit()
	}

	return nil
}

// verifyCommit verifies if the received commit message is equivalent to our subject
func (c *core) verifyD_Commit(commit *poDC.Subject, src poDC.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptD_Commit(msg *message, src poDC.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the commit message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	return nil
}
