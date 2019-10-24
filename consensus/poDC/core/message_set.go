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
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/consensus/poDC"  //podc
)

// Construct a new message set to accumulate messages for given sequence/view number.
func newMessageSet(valSet poDC.ValidatorSet) *messageSet {
	return &messageSet{
		view: &poDC.View{  // yichoi for podc
			Round:    new(big.Int),
			Sequence: new(big.Int),
		},
		messagesMu: new(sync.Mutex),
		messages:   make(map[common.Hash]*message),
		valSet:     valSet,
	}
}

// ----------------------------------------------------------------------------

type messageSet struct {
	view       *poDC.View
	valSet     poDC.ValidatorSet
	messagesMu *sync.Mutex
	messages   map[common.Hash]*message
}
// yichoi begin
type messageSet struct {
	view       *poDC.View
	valSet     poDC.ValidatorSet
	messagesMu *sync.Mutex
	messages   map[common.Hash]*message
}
// end

func (ms *messageSet) View() *poDC.View {  //yichoi
	return ms.view
}

func (ms *messageSet) Add(msg *message) error {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	if err := ms.verify(msg); err != nil {
		return err
	}

	return ms.addVerifiedMessage(msg)
}

func (ms *messageSet) Values() (result []*message) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	for _, v := range ms.messages {
		result = append(result, v)
	}

	return result
}

func (ms *messageSet) Size() int {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return len(ms.messages)
}

// ----------------------------------------------------------------------------

func (ms *messageSet) verify(msg *message) error {
	// verify if the message comes from one of the validators
	if _, v := ms.valSet.GetByAddress(msg.Address); v == nil {
		return poDC.ErrUnauthorizedAddress
	}

	// TODO: check view number and sequence number

	return nil
}

func (ms *messageSet) addVerifiedMessage(msg *message) error {
	ms.messages[poDC.RLPHash(msg)] = msg
	return nil
}

func (ms *messageSet) String() string {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	addresses := make([]string, 0, len(ms.messages))
	for _, v := range ms.messages {
		addresses = append(addresses, v.Address.String())
	}
	return fmt.Sprintf("[%v]", strings.Join(addresses, ", "))
}
