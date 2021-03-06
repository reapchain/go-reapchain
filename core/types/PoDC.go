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

package types

import (
"errors"
"io"

"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/rlp"
)

var (
	// IstanbulDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine
	PoDCDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	PoDCExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	PoDCExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidPoDCHeaderExtra = errors.New("invalid PoDC header extra-data")
)


// Reapchain PoDC extra struct
type PoDCExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *PoDCExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.Validators,
		ist.Seal,
		ist.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder, and load the istanbul fields from a RLP stream.
func (ist *PoDCExtra) DecodeRLP(s *rlp.Stream) error {
	var PoDCExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := s.Decode(&PoDCExtra); err != nil {
		return err
	}

	ist.Validators, ist.Seal, ist.CommittedSeal = PoDCExtra.Validators, PoDCExtra.Seal, PoDCExtra.CommittedSeal

	return nil
}

// ExtractIstanbulExtra extracts all values of the IstanbulExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
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


// PoDCFilteredHeader returns a filtered header which some information (like seal, committed seals)
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

