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
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/core/types"
	elog "github.com/ethereum/go-ethereum/log"
)

func makeBlock(number int64) *types.Block {
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(number),
		GasLimit:   big.NewInt(0),
		GasUsed:    big.NewInt(0),
		Time:       big.NewInt(0),
	}
	block := &types.Block{}
	return block.WithSeal(header)
}

func newTestProposal() podc.Proposal {
	return makeBlock(1)
}

func TestNewRequest(t *testing.T) {
	testLogger.SetHandler(elog.StdoutHandler)

	N := uint64(28)  //일반 Validators
	F := uint64(1)  //배신자 노드 한개 가정
	elog.Info("current params.MaximumExtraDataSize=%d", len(int(params.MaximumExtraDataSize))) //check for Byte or Kilo Byte?
	sys := NewTestSystemWithBackend(N, F)

	close := sys.Run(true)
	defer close()

	request1 := makeBlock(1)
	sys.backends[0].NewRequest(request1)

	select {
	case <-time.After(1 * time.Second):
	}

	request2 := makeBlock(2)
	sys.backends[0].NewRequest(request2)

	select {
	case <-time.After(1 * time.Second):
	}

	for _, backend := range sys.backends {
		if len(backend.commitMsgs) != 2 {
			t.Errorf("the number of executed requests mismatch: have %v, want 2", len(backend.commitMsgs))
		}
		if !reflect.DeepEqual(request1.Number(), backend.commitMsgs[0].Number()) {
			t.Errorf("the number of requests mismatch: have %v, want %v", request1.Number(), backend.commitMsgs[0].Number())
		}
		if !reflect.DeepEqual(request2.Number(), backend.commitMsgs[1].Number()) {
			t.Errorf("the number of requests mismatch: have %v, want %v", request2.Number(), backend.commitMsgs[1].Number())
		}
	}
}
