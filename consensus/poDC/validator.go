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
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// 모든 Validator 로 선언하고, 여기서 각 Label/ tag 붙여서 관리한다.
type Validator interface {
	// Address returns address
	Address() common.Address

	// String representation of Validator
	String() string

	Tag() uint64 // Tag()  is bug ,, just test by yichoi
	// Tag define : Senator(상원), parliamentarian(하원), general(일반), candidate(운영위 후보)

	Qrnd() uint64
}

// ----------------------------------------------------------------------------

type Validators []Validator // go 배열 표현

/* 설명 : Validators = [ address, String, Tag ] [ ... ] [ ... ] ......... */
/*  Validator node
|-------------------------|
|  enode address(20 byte)
|-------------------------|
   String인데,,String()모르겠으나, 그냥 문자열,
   Validator 를 나타내는
|-------------------------|
   Tag를 새로 달았음.
|-------------------------|        ............   N 개 Validators



*/
func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------------------------------------

type ValidatorSet interface {
	// Calculate the proposer
	CalcProposer(lastProposer common.Address, round uint64) // 최초 Proposer 가 누군지 계산 ,,

	// 코디 후보군 선정 ?
	RecvCordinator(lastProposer common.Address, round uint64) //yichoi

	// Return the validator size
	Size() int
	// Return the validator array
	List() []Validator
	// Get validator by index
	GetByIndex(i uint64) Validator
	// Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
	// Get current proposer
	GetProposer() Validator

	// 상임위, 운영위 확정 구성 위원회 정보 가져오기 , Steering committee ( 운영위원회 )
	GetConfirmedCommittee() Validator //yichoi

	// Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// Is request a ExtraDATA to Qmanager
	IsRequestQman(address common.Address) bool //podc

	// Add validator
	AddValidator(address common.Address) bool
	// Remove validator
	RemoveValidator(address common.Address) bool
	// Copy validator set
	Copy() ValidatorSet
	// Get the maximum number of faulty nodes
	F() int
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator
