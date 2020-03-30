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
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rlp"
)
//type ValidatorInfo struct {
//	Address common.Address
//	Tag istanbul.Tag
//	Qrnd uint64
//}

// type ValidatorInfos []ValidatorInfo

type Engine interface {

	//Start(lastSequence *big.Int, lastProposer common.Address, lastProposal istanbul.Proposal, qmanager []*discover.Node) error  //modified by yichoi for qmanager
	Start(lastSequence *big.Int, lastProposer common.Address, lastProposal podc.Proposal, qmanager []*discover.Node) error

	Stop() error
}

type State uint64



const (

	StateAcceptRequest State = iota
	StatePreprepared  // = pre-prepare of podc
	StatePrepared     // = d-select
	StateCommitted    // = d-commit
	StateRequestQman  // request to Qman to get ExtraData

/*
	StatePreprepared
	StatePrepared
	StateCommitted
	StateRequestQMan  // state to send request for extra data to Qmanager
	StateAcceptQMan  //podc
	StateD_selected  //podc
	StateD_committed //podc
>>>>>>> working:consensus/podc/core/types.go */
)

func (s State) String() string {
	if s == StateRequestQman {
		return "Request ExtraData"   //새 라운드가 시작하면, 매번 Qmanager에게 ExtraDATA 요청한다.
	} else if s == StateAcceptRequest {
		return "Accept request"
	} else if s == StatePreprepared {
		return "Preprepared"
	} else if s == StatePrepared {
		return "Prepared"
	} else if s == StateCommitted {
		return "Committed"
	} else {
		return "Unknown"
	}
	/*  else if s == StateD_selected { //podc
		return "D_selected"
	} else if s == StateD_committed { //podc
		return "D_committed"
	} */

}

// Cmp compares s and y and returns:
//   -1 if s is the previous state of y
//    0 if s and y are the same state
//   +1 if s is the next state of y
func (s State) Cmp(y State) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}

const (
	msgPreprepare uint64 = iota
	msgDSelect
	msgCoordinatorDecide
	msgRacing
	msgCandidateDecide
	msgPrepare
	msgCommit
	msgRoundChange

	// msgD_commit          // podc
	// msgD_select          // podc
	msgRequestQman
	msgHandleQman // for Qman, receive event handler from geth ( sending qmanager )
	msgAll
	//New Qmanager Implementation
	msgExtraDataRequest
	msgExtraDataSend

	msgCoordinatorConfirmRequest
	msgCoordinatorConfirmSend

//=======
/*	msgRequestQman        // send a request to Qmanager in order to get ExtraDATA
	msgReceivedFromQman   // When receive a message from Qmanager, msg is ExtraData
	msgPrepare
	msgCommit
	msgRoundChange

	msgGetCandiateList   // podc
	msgStartRacing       // podc
	msgRegisterCommittee // podc

//>>>>>>> working:consensus/podc/core/types.go
*/
)

type message struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *message) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}

	// Validate message (on a message without Signature)
	if validateFn != nil {
		var payload []byte
		payload, err = m.PayloadNoSig()
		if err != nil {
			return err
		}

		_, err = validateFn(payload, m.Signature)
	}
	// Still return the message even the err is not nil
	return err
}

func (m *message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&message{
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// 메시지를 주면 val 값으로 돌려줌
func (m *message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val) //DecodeBytes parses RLP data from b into val.
}

func (m *message) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v}", m.Code, m.Address.String())
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

// ----------------------------------------------------------------------------

type roundChange struct {
	Round    *big.Int
	Sequence *big.Int
	Digest   common.Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (rc *roundChange) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		rc.Round,
		rc.Sequence,
		rc.Digest,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (rc *roundChange) DecodeRLP(s *rlp.Stream) error {
	var rawRoundChange struct {
		Round    *big.Int
		Sequence *big.Int
		Digest   common.Hash
	}
	if err := s.Decode(&rawRoundChange); err != nil {
		return err
	}
	rc.Round, rc.Sequence, rc.Digest = rawRoundChange.Round, rawRoundChange.Sequence, rawRoundChange.Digest
	return nil
}

// ----------------------------------------------------------------------------
/*
type extraData struct {
	HashValue common.Hash
	Details   details
	Signature []byte
}

func (ed *extraData) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ed.HashValue,
		ed.Details,
		ed.Signature,
	})
}

func (ed *extraData) DecodeRLP(s *rlp.Stream) error {
	var extraData struct {
		HashValue common.Hash
		Details   Validator
		Signature []byte
	}
	if err := s.Decode(&extraData); err != nil {
		return err
	}
	ed.HashValue, ed.Details, ed.Signature = extraData.HashValue, extraData.Details, extraData.Signature
	return nil
} */
