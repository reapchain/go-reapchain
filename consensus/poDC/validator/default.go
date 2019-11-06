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

package validator

import (
	"math"
	"reflect"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/poDC"

	"github.com/ethereum/go-ethereum/common"
)

const (
	Senator         uint64 = iota // 상원
	Parliamentarian               // 하원
	Candidate                     // 하원, 운영위 후보군
	General                       // 일반 노드, 상원, 하원도 아닌.
	QManager                      // Q-Manager
)

type defaultValidator struct {
	address common.Address
	tag     uint64
}

func (val *defaultValidator) Address() common.Address {
	return val.address
}

func (val *defaultValidator) String() string {
	return val.Address().String()
}

func (val *defaultValidator) Tag() uint64 {
	return val.tag
}

// ----------------------------------------------------------------------------

/* yichoi unmarked for replace istanbul with PoDC default validator set
type defaultSet struct {
	validators  poDC.Validators  // 상임위, 운영위, 운영위 후보, 코디로 나누는 방법 연구
	proposer    poDC.Validator  // front node for podc
	validatorMu sync.RWMutex

	selector poDC.ProposalSelector
} */
type defaultSet struct {
	validators poDC.Validators // 상임위, 운영위, 운영위 후보, 코디로 나누는 방법 연구

	quantumManager poDC.Validator // Quantum manager, 일반노드 같이 enode 번호를 갖게끔 설계
	coordinator    poDC.Validator //

	proposer    poDC.Validator // = front node for podc
	validatorMu sync.RWMutex

	selector poDC.ProposalSelector
}

// 최초 Proposer 는 newDefaultSet를 호출할때, sort 후 GetByIndex (0) 로 선택하고, 추후에는 selector를 사용해서 한다.

func newDefaultSet(addrs []common.Address, selector poDC.ProposalSelector) *defaultSet {
	valSet := &defaultSet{}

	// init validators
	valSet.validators = make([]poDC.Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = New(addr)
		// 추후 Validators 를 로그로 찍어볼것
	}
	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer = valSet.GetByIndex(0)
	}
	//set proposal selector : front node
	valSet.selector = selector

	//get coordinator and steering committee from Qmanager
	//get coordinator
	//valSet.coordinator = func (valSet *defaultSet) RecvCoordinator(lastProposer common.Address, round uint64)

	//get candidates  of steering committe
	//  valSet.GetConfirmedCommittee() = .......determin and select steering committee

	return valSet
}

func (valSet *defaultSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

func (valSet *defaultSet) List() []poDC.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

func (valSet *defaultSet) GetByIndex(i uint64) poDC.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	if i < uint64(valSet.Size()) {
		return valSet.validators[i]
	}
	return nil
}

func (valSet *defaultSet) GetByAddress(addr common.Address) (int, poDC.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Address() {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetProposer() poDC.Validator {
	return valSet.proposer
}

// Check that " I'm coordinator or not " , self check.
// 코디네이트는 핸들어에서 브로드캐스트 메시지를 받으면,
//
func (valSet *defaultSet) SelfCheckCoordi() poDC.Validator {
	return valSet.coordinator //"
}

//yichoi coordinator infomation is received from Qmanger server
func (valSet *defaultSet) RecvCoordinator(lastProposer common.Address, round uint64) { //?
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	// Quantum manager로부터 Coordinator 정보를 가져온다.

	valSet.coordinator = valSet.selector(valSet, lastProposer, round) //Front node 선택 계산.
}

//
// Check that "I'm a candidate of steering committee from Qmanager"

func (valSet *defaultSet) SelfCheckCandidate() poDC.Validator {
	return valSet.candidate //
}

//

//yichoi Get Confirmed Committee
func (valSet *defaultSet) GetConfirmedCommittee() poDC.Validator {
	return valSet.proposer
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	valSet.proposer = valSet.selector(valSet, lastProposer, round) //Front node 선택 계산.
}

func calcSeed(valSet poDC.ValidatorSet, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}

func emptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func roundRobinProposer(valSet poDC.ValidatorSet, proposer common.Address, round uint64) poDC.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round) + 1
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func stickyProposer(valSet poDC.ValidatorSet, proposer common.Address, round uint64) poDC.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round)
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

//Quantum Random Finding .. ?
func qrfProposer(valSet poDC.ValidatorSet, proposer common.Address, round uint64) poDC.Validator {

	/* Quantum manager 와 주고 받는 것 구현 ?
	    // ExtraDATA에서,, 코디와 운영위, 상임위등, Validator들의 정보를 파싱해서, 수집


		로부터 데이타 수신해서,
		Extra data를 수신한다.
	*/

}

func (valSet *defaultSet) AddValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for _, v := range valSet.validators {
		if v.Address() == address {
			return false
		}
	}
	valSet.validators = append(valSet.validators, New(address))
	// TODO: we may not need to re-sort it again
	// sort validator
	sort.Sort(valSet.validators)
	return true
}

func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.Address() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

func (valSet *defaultSet) Copy() poDC.ValidatorSet {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	addresses := make([]common.Address, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		addresses = append(addresses, v.Address())
	}
	return newDefaultSet(addresses, valSet.selector)
}

func (valSet *defaultSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }
