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

	StateRequest State = iota // request to Qman to get ExtraData
	StateAcceptRequest // geth가 프로포저(프런트노드)가 아닌 경우 일반 노드들은 수동적으로 외부 노드로부터 메시지를 수신 대기 상태 유지
	StatePreprepared  // = pre-prepare of podc : Bx including coordi info, broadcast to all nodes.
	StateDselected     // = d-select
	StateDCommitted    // = d-commit
    StateFinalCommitted     // variables initialize before go New Round state


	//
	//StateAcceptRequest State = iota
	//StatePreprepared  // = pre-prepare of podc
	//StatePrepared     // = d-select
	//StateCommitted    // = d-commit
	//StateRequestQman  // request to Qman to get ExtraData
)

func (s State) String() string {
	if s == StateRequest {
		return "Request ExtraData to Qman"   //새 라운드가 시작하면, 매번 Qmanager에게 ExtraDATA 요청한다.
	}  else if s == StatePreprepared {
		return "Preprepared"
	} else if s == StateDselected {
		return "Dselected"
	} else if s == StateDCommitted {
		return "DCommitted"
	}else if s == StateFinalCommitted {
		return "FinalCommitted"
	} else {
		return "Unknown"
	}
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

	//1. request step
	msgRequest  uint64 = iota  //1.  	msgHandleQman // for Qman, receive event handler from geth ( sending qmanager )

	//2. Pre-prepare step
	msgPreprepare
	msgPrepare
	//3. D-select step
	msgDSelect
	msgCoordinatorDecide  //notify to Qman
	msgRacing             // between coordi and candidates
	msgCandidateDecide    //notify to Coordi

	//3. D-commit step
	msgDCommit

	msgExtraDataRequest  //Request to Qman
	msgExtraDataSend     //Qman send to geth

	msgCoordinatorConfirmRequest  //Geth send Qman
	msgCoordinatorConfirmSend     //Qman sennd geth

	// etc:
	msgRoundChange
	/* In paper:
	1. Request : 프런트 노드는 Qmanager에게 접속, 코디 정보가 담긴 ExtraData를 가져오는 단계
	2. Pre-prepare : 프런트 노드는 코디 정보가 담긴 ExtraData로 Bx를 생성 후 모든 노드에 브로드캐스트 하는 단계
	3. D-select :
	   코디 : 본인의 패킷을 까서, 자신이 코디이면, Qmanager에게 코디임을 등록 하고, 레이싱 메시지를 각 운영위 후보군에 멀티캐스트
	         코디는 메시지 이벤트 핸들러에서, 운영위 후보군으로부터, 레이싱 메시지를 선착순으로 선발해서, 15등 안에 도달하는 운영위 후보만,
	         선발하여, 최종 운영위 후보( 하원 )로 확정하고, Qmanager에 확정 메시지와, 선발된 운영위 후보 목록을 전송한다.
	   운영위 후보군: 본인의 패킷을 까서, 자신이 운영위 후보군이면 코디로 부터 레이싱메시지 수신을 기다렸다가, 레이싱 메신지를 수신하면,
	         코디에게 "내가 운영위 후보"라는 메시지를 보내고 레이싱에 참여한다.
	4. D-commit : 코디는 , 상임위, 최종 선발된 운영위들에 대하여 , 제안된 블럭을 가지고, , ( 상임위 14 + 확정된 운영위(하원) 15= 총 29개)
	              투표(Voting)를 통하여 51%의 동의 메시지를 수신하면, 전체 노드에게 결과를 브로드캐스트해서, 각 노드들이 블럭을 체인에 삽입하도록 명령한다.
	5. D-committed: 블럭 삽입이 성공한 후,  이 상태로 바뀐다.
	6. Final committed: New Round 로 가기 전 단계로, 모든 변수 등을 초기화 한다.
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
