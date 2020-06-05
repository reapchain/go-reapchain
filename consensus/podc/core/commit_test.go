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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc/validator"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestHandleCommit(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

	proposal := newTestProposal()
	expectedSubject := &podc.Subject{
		View: &podc.View{
			Round:    big.NewInt(0),
			Sequence: proposal.Number(),
		},
		Digest: proposal.Hash(),
	}

	testCases := []struct {
		system      *testSystem
		expectedErr error
	}{
		{
			// normal case
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&podc.View{
							Round:    big.NewInt(0),
							Sequence: big.NewInt(1),
						},
						c.valSet,
					)

					if i == 0 {
						// replica 0 is primary
						c.state = StateDselected
					}
				}
				return sys
			}(),
			nil,
		},
		{
			// future message
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is primary
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StatePreprepared
					} else {
						c.current = newTestRoundState(
							&podc.View{
								Round:    big.NewInt(2),
								Sequence: big.NewInt(3),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errFutureMessage,
		},
		{
			// subject not match
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is primary
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StatePreprepared
					} else {
						c.current = newTestRoundState(
							&podc.View{
								Round:    big.NewInt(0),
								Sequence: big.NewInt(0),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errOldMessage,
		},
		{
			// jump state
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&podc.View{
							Round:    big.NewInt(0),
							Sequence: proposal.Number(),
						},
						c.valSet,
					)

					// only replica0 stays at StatePreprepared
					// other replicas are at StatePrepared
					if i != 0 {
						c.state = StateDselected
					} else {
						c.state = StatePreprepared
					}
				}
				return sys
			}(),
			nil,
		},
		// TODO: double send message
	}

OUTER:
	for _, test := range testCases {
		test.system.Run(false)

		v0 := test.system.backends[0]
		r0 := v0.engine.(*core)

		for i, v := range test.system.backends {
			validator := r0.valSet.GetByIndex(uint64(i))
			m, _ := Encode(v.engine.(*core).current.Subject())
			if err := r0.handleDCommit(&message{
				Code:          msgDCommit,
				Msg:           m,
				Address:       validator.Address(),
				Signature:     []byte{},
				CommittedSeal: validator.Address().Bytes(), // small hack
			}, validator); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
				continue OUTER
			}
		}

		// prepared is normal case
		if r0.state != StateDCommitted {
			// There are not enough commit messages in core
			if r0.state != StateDselected {
				t.Errorf("state mismatch: have %v, want %v", r0.state, StateDselected)
			}
			if r0.current.Dcommits.Size() > 2*r0.valSet.F() {
				t.Errorf("the size of commit messages should be less than %v", 2*r0.valSet.F()+1)
			}

			continue
		}

		// core should have 2F+1 prepare messages
		if r0.current.Dcommits.Size() <= 2*r0.valSet.F() {
			t.Errorf("the size of commit messages should be larger than 2F+1: size %v", r0.current.Dcommits.Size())
		}

		// check signatures large than 2F+1
		signedCount := 0
		signers := make([]common.Address, len(v0.committedSeals[0])/common.AddressLength)
		for i := 0; i < len(signers); i++ {
			copy(signers[i][:], v0.committedSeals[0][i*common.AddressLength:])
		}
		for _, validator := range r0.valSet.List() {
			for _, signer := range signers {
				if validator.Address() == signer {
					signedCount++
					break
				}
			}
		}
		if signedCount <= 2*r0.valSet.F() {
			t.Errorf("the expected signed count should be larger than %v, but got %v", 2*r0.valSet.F(), signedCount)
		}
	}
}

// round is not checked for now
func TestVerifyCommit(t *testing.T) {
	// for log purpose
	privateKey, _ := crypto.GenerateKey()
	peer := validator.New(getPublicKeyAddress(privateKey))
	valSet := validator.NewSet([]common.Address{peer.Address()}, podc.RoundRobin)

	sys := NewTestSystemWithBackend(uint64(1), uint64(0))

	testCases := []struct {
		expected   error
		commit     *podc.Subject
		roundState *roundState
	}{
		{
			// normal case
			expected: nil,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// old message  1.
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// different digest  2.
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: common.StringToHash("1234567890"),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// malicious package(lack of sequence)  3. got="{View: %!v(PANIC=String method: runtime error: invalid memory address or nil pointer dereference),
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(0), Sequence: nil},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// wrong prepare message with same sequence but different round
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(1), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// wrong prepare message with same round but different sequence
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(0), Sequence: big.NewInt(1)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
//begin - by yichoi
		{
			// wrong prepare message with same round but different sequence
			expected: errInconsistentSubject,
			commit: &podc.Subject{
				View:   &podc.View{Round: big.NewInt(1), Sequence: nil }, //got="{View: {Round: 1, Sequence: 1},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&podc.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},  //expected="{View: {Round: 2, Sequence: 2},
				valSet,
			),
		},
//end
	}
	for i, test := range testCases {
		c := sys.backends[0].engine.(*core)
		c.current = test.roundState

		//if err := c.verifyCommit(test.commit, peer); err != nil {
		//	if err != test.expected {
		//		t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
		//	}
		//}
		if err := c.verifyDCommit(test.commit, peer); err != nil {
			if err != test.expected {
				t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
			}
		}
	}
}
