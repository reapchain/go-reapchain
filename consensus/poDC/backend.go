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

package poDC

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {
	// Address returns the owner's address
	Address() common.Address

	// Validators returns the validator set
	Validators(proposal Proposal) ValidatorSet  // old istanbul method, I'll remove it in the future after debuging

	// Validators returns the validator set
	EmptyValidators(proposal Proposal) string   // PoDC 상원, 하원,


	// EventMux returns the event mux in backend
	EventMux() *event.TypeMux

	// Send sends a message to specific target : 특정 노드에 보낼때,
	Send(payload []byte, target common.Address) error

	// Broadcast sends a message to all validators : 전체 노드에 보낼때, 기 정의된 Validator집합에..
	Broadcast(valSet ValidatorSet, payload []byte) error

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposal Proposal, seals []byte) error

	// NextRound is called when we want to trigger next Seal()
	NextRound() error

	// Verify verifies the proposal.
	Verify(Proposal) error

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	// CheckSignature verifies the signature by checking if it's signed by
	// the given validator
	CheckSignature(data []byte, addr common.Address, sig []byte) error
}

//EmptyValidators(proposal Proposal) bool