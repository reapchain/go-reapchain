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

// made by yichoi for PoDC
// PoDC ExtraDATA 타입 정의와 헤더 등이 맞는지 검증 하는 모듈
package types

import (
"errors"
"io"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/rlp"
)
//yichoi
var (
	// IstanbulDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine
	PoDCDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")
    // 블럭이 이스탄불 합의 엔진으로부터 온건지 구분하는 다이제스트 ?

	PoDCExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	PoDCExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidPoDCHeaderExtra = errors.New("invalid PoDC header extra-data")
)


// Reapchain PoDC extra struct
type PoDCExtra struct {
	Validators    []common.Address  //다른가? podc에서는 ?  그냥 20바이트 enode 주소 구조체 아님. 20 바이트 배열 한개
	//Validator들이 3종류, 상원, 하원 운영위, 운영위 후보, 일반 검증자,,로 세분화됨...
	// 이걸 나누는 것보다,  Tag를 다는게 효율적일듯.
	Seal          []byte
	CommittedSeal [][]byte
	// Tag           []byte   // podc validator 종류 표시 태크  , 1, 상원, 2, 하원 3, 운영위 후보, Candidates, 4. 일반노드
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *PoDCExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.Validators,
		ist.Seal,
		ist.CommittedSeal,
		//ist.Tag, //podc
	})
}

// DecodeRLP implements rlp.Decoder, and load the istanbul fields from a RLP stream.
func (ist *PoDCExtra) DecodeRLP(s *rlp.Stream) error {
	var PoDCExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
		// Tag           []byte //podc
	}
	if err := s.Decode(&PoDCExtra); err != nil {
		return err
	}
	// disableed : ist.Validators, ist.Seal, ist.CommittedSeal, ist.Tag = PoDCExtra.Validators, PoDCExtra.Seal, PoDCExtra.CommittedSeal, PoDCExtra.Tag  //by yichoi old , incompleted
	ist.Validators, ist.Seal, ist.CommittedSeal = PoDCExtra.Validators, PoDCExtra.Seal, PoDCExtra.CommittedSeal //Taehun implemetation
	// ?  PoDCExtra.Tag   추가 ?
	return nil
}

// ExtractIstanbulExtra extracts all values of the IstanbulExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.

// made by yichoi for podc ExtractPoDCExtra data
// 블럭헤더 정보의 ExtraData를 가져오는데, 아래 함수를 Qmanager 서버로부터 가져오는 것으로 바꿔야함.
func ExtractPoDCExtra(h *Header) (*PoDCExtra, error) {
	if len(h.Extra) < PoDCExtraVanity {
		return nil, ErrInvalidPoDCHeaderExtra
	}

	var PoDCExtra *PoDCExtra
	err := rlp.DecodeBytes(h.Extra[PoDCExtraVanity:], &PoDCExtra)
	if err != nil {
		return nil, err
	}
	return PoDCExtra, nil
}

/* func RequestPoDCExtraToQman(Qmanager common.Address) ( error) {


To do future

	/* send request to Qmanager

	하고 에러 없는지만 체크




	handler 함수에서,, Qmanger로부터 ,, ExtraData를 받는 부분 구현.

	요청 메시지 보내고,, state는 Extradata waiting 상태로 변경, */




//return error

// }

// IstanbulFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Istanbul hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.


func PoDCFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	PoDCExtra, err := ExtractPoDCExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		PoDCExtra.Seal = []byte{}
	}
	PoDCExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&PoDCExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:PoDCExtraVanity], payload...)

	return newHeader
}

