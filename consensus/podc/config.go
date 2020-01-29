
// Copyright 2019 Reapchain
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

package podc

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota   //?
	Sticky
	QRF   //yichoi 립체인 퀀텀메카니즘 합의 정책 상수, 1단계 : 라운드로빈 방식으로 ,, 2단계에서 Quntum Random 으로
)

type Config struct {
	RequestTimeout uint64         `toml:",omitempty"` // The timeout for each Istanbul round in milliseconds. This timeout should be larger than BlockPauseTime
	BlockPeriod    uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	BlockPauseTime uint64         `toml:",omitempty"` // Delay time if no tx in block, the value should be larger than BlockPeriod
	ProposerPolicy ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	Epoch          uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
}

var DefaultConfig = &Config{
	RequestTimeout: 10000,
	BlockPeriod:    1,
	BlockPauseTime: 2,
	ProposerPolicy: RoundRobin,   // pilot 으로 1차는 proposer 선택 정책은 라운드로빈 유지, 2차로 수정 예정,
	                               // 1차 디버깅 후, QRF로 바꿀것,
	Epoch:          30000,
}

