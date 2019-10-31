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
)
// Pre-Prepare state 에서 Qmanager에서 ExtraData를 수신하는 행위도 ConsensusDataEvent로 분류할 것.
type ConsensusDataEvent struct {
	// target to send message
	Target common.Address
	// consensus message data
	Data []byte
}
// RequestEvent 는 Qmanager와 주고 받는 것을 포함해야 코드가 심플할 듯.
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
