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
	"io"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/rlp"
)

func newRoundState(view *podc.View, validatorSet podc.ValidatorSet) *roundState {
	return &roundState{
		round:       view.Round,
		sequence:    view.Sequence,
		Preprepare:  nil,
		Prepares:    newMessageSet(validatorSet),
		Commits:     newMessageSet(validatorSet),
		Checkpoints: newMessageSet(validatorSet),
		mu:          new(sync.RWMutex),
	}
}

// roundState stores the consensus state
type roundState struct {
	round       *big.Int
	sequence    *big.Int
	Preprepare  *istanbul.Preprepare
	Prepares    *messageSet
	Commits     *messageSet
	Checkpoints *messageSet

	mu *sync.RWMutex
}

func (s *roundState) Subject() *podc.Subject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare == nil {
		return nil
	}

	return &podc.Subject{
		View: &podc.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest: s.Preprepare.Proposal.Hash(),
	}
}

func (s *roundState) SetPreprepare(preprepare *podc.Preprepare) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Preprepare = preprepare
}
/* begin : yichoi added for d-select and d-commit */
/* func (s *roundState) SetD_select(d_select *podc.Preprepare) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.D_select = d_select
}

func (s *roundState) SetD_commit(d_commit *podc.Preprepare) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.D_commit = d_commit
} */
/* end */
//왜 락을 걸지 ?
func (s *roundState) Proposal() podc.Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare != nil {
		return s.Preprepare.Proposal  //제안할 블럭을 가져온다... 합의가 끝나면,, 체인에 연결할 블럭을 가져온다.
	}

	return nil
}

func (s *roundState) SetRound(r *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.round = new(big.Int).Set(r)
}

func (s *roundState) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState) SetSequence(seq *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence = seq
}

func (s *roundState) Sequence() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sequence
}

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var ss struct {
		Round       *big.Int
		Sequence    *big.Int
		Preprepare  *podc.Preprepare
		Prepares    *messageSet
		Commits     *messageSet
		Checkpoints *messageSet
	}

	if err := stream.Decode(&ss); err != nil {
		return err
	}
	s.round = ss.Round
	s.sequence = ss.Sequence
	s.Preprepare = ss.Preprepare
	s.Prepares = ss.Prepares
	s.Commits = ss.Commits
	s.Checkpoints = ss.Checkpoints
	s.mu = new(sync.RWMutex)

	return nil
}

// EncodeRLP should write the RLP encoding of its receiver to w.
// If the implementation is a pointer method, it may also be
// called for nil pointers.
//
// Implementations should generate valid RLP. The data written is
// not verified at the moment, but a future version might. It is
// recommended to write only a single value but writing multiple
// values or no value at all is also permitted.
func (s *roundState) EncodeRLP(w io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return rlp.Encode(w, []interface{}{
		s.round,
		s.sequence,
		s.Preprepare,
		s.Prepares,
		s.Commits,
		s.Checkpoints,
	})
}
