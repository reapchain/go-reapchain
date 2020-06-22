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

package podc

import (
	"github.com/ethereum/go-ethereum/common"
)

type ConsensusDataEvent struct {
	// target to send message
	Target common.Address
	// consensus message data
	Data []byte
}

type QmanDataEvent struct {
	// target to send message
	Target common.Address
	// consensus message data
	Data []byte
}

type RequestEventQman struct {
	Proposal Proposal
}
//-----------------------------------------
type RequestEvent struct {
	Proposal Proposal
}

type MessageEvent struct {
	Payload []byte
}

type FinalCommittedEvent struct {
	Proposal Proposal
	Proposer common.Address
}
