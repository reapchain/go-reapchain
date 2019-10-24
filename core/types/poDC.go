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

	PoDCExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	PoDCExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidPoDCHeaderExtra = errors.New("invalid poDC header extra-data")
)
/* unmarked by yichoi
type IstanbulExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
} */

// Reapchain PoDC extra struct
type PoDCExtra struct {
	Validators    []common.Address  //다른가? podc에서는 ?
	//Validator들이 3종류, 상원, 하원 운영위, 운영위 후보, 일반 검증자,,로 세분화됨...
	// 이걸 나누는 것보다,  Tag를 다는게 효율적일듯.
	Seal          []byte
	CommittedSeal [][]byte
	Tag           []byte   // podc validator 종류 표시 태크  , 1, 상원, 2, 하원 3, 운영위 후보, Candidates, 4. 일반노드
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *PoDCExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.Validators,
		ist.Seal,
		ist.CommittedSeal,
		ist.Tag, //podc
	})
}

// DecodeRLP implements rlp.Decoder, and load the istanbul fields from a RLP stream.
func (ist *PoDCExtra) DecodeRLP(s *rlp.Stream) error {
	var poDCExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
		Tag           []byte //podc
	}
	if err := s.Decode(&poDCExtra); err != nil {
		return err
	}
	ist.Validators, ist.Seal, ist.CommittedSeal = poDCExtra.Validators, poDCExtra.Seal, poDCExtra.CommittedSeal
	// ?  poDCExtra.Tag   추가 ?
	return nil
}

// ExtractIstanbulExtra extracts all values of the IstanbulExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.

// made by yichoi for podc ExtractPoDCExtra data
func ExtractPoDCExtra(h *Header) (*PoDCExtra, error) {
	if len(h.Extra) < PoDCExtraVanity {
		return nil, ErrInvalidPoDCHeaderExtra
	}

	var poDCExtra *PoDCExtra
	err := rlp.DecodeBytes(h.Extra[PoDCExtraVanity:], &poDCExtra)
	if err != nil {
		return nil, err
	}
	return poDCExtra, nil
}

// IstanbulFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Istanbul hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.


func PoDCFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	poDCExtra, err := ExtractPoDCExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		poDCExtra.Seal = []byte{}
	}
	poDCExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(&poDCExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:PoDCExtraVanity], payload...)

	return newHeader
}

