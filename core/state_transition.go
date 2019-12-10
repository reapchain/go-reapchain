// Copyright 2014 The go-ethereum Authors
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
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var (
	Big0                         = big.NewInt(0)
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *GasPool
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas *big.Int
	value      *big.Int
	data       []byte
	state      vm.StateDB

	evm *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	//FromFrontier() (common.Address, error)
	To() *common.Address

	GasPrice() *big.Int
	Gas() *big.Int
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte

	Governance() bool   // yhheo
}

// IntrinsicGas computes the 'intrinsic gas' for a message
// with the given data.
//
// TODO convert to uint64
func IntrinsicGas(data []byte, contractCreation, homestead bool) *big.Int {
	igas := new(big.Int)
	if contractCreation && homestead {
		igas.SetUint64(params.TxGasContractCreation)
	} else {
		//igas.SetUint64(params.TxGas)
		igas.SetUint64(0) // KBW
	}

	log.Info("IntrinsicGas","gas1", igas, "data length", len(data)) // KBW

	if len(data) > 0 {
		var nz int64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		m := big.NewInt(nz)
		m.Mul(m, new(big.Int).SetUint64(params.TxDataNonZeroGas))
		igas.Add(igas, m)
		m.SetInt64(int64(len(data)) - nz)
		m.Mul(m, new(big.Int).SetUint64(params.TxDataZeroGas))
		igas.Add(igas, m)
	}
	return igas
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:         gp,
		evm:        evm,
		msg:        msg,
		gasPrice:   msg.GasPrice(),
		initialGas: new(big.Int),
		value:      msg.Value(),
		data:       msg.Data(),
		state:      evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool) ([]byte, *big.Int, error) {
	st := NewStateTransition(evm, msg, gp)

	ret, _, gasUsed, err := st.TransitionDb()
	return ret, gasUsed, err
}

func (st *StateTransition) from() vm.AccountRef {
	f := st.msg.From()
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return vm.AccountRef(f)
}

func (st *StateTransition) to() vm.AccountRef {
	if st.msg == nil {
		return vm.AccountRef{}
	}
	to := st.msg.To()
	if to == nil {
		return vm.AccountRef{} // contract creation
	}

	reference := vm.AccountRef(*to)
	if !st.state.Exist(*to) {
		st.state.CreateAccount(*to)
	}
	return reference
}

func (st *StateTransition) useGas(amount uint64) error {
	log.Info("TransitionDb","st.gas", st.gas,"amount", amount) // KBW

	if st.gas < amount {
		return vm.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	mgas := st.msg.Gas()
	if mgas.BitLen() > 64 {
		return vm.ErrOutOfGas
	}

	mgval := new(big.Int).Mul(mgas, st.gasPrice)

	var (
		state  = st.state
		sender = st.from()
	)
	if state.GetBalance(sender.Address()).Cmp(mgval) < 0 {
		return errInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(mgas); err != nil {
		return err
	}
	st.gas += mgas.Uint64()

	st.initialGas.Set(mgas)
	state.SubBalance(sender.Address(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	msg := st.msg
	sender := st.from()

	// Make sure this transaction's nonce is correct
	if msg.CheckNonce() {
		if n := st.state.GetNonce(sender.Address()); n != msg.Nonce() {
			return fmt.Errorf("invalid nonce: have %d, expected %d", msg.Nonce(), n)
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and returning the result
// including the required gas for the operation as well as the used gas. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, requiredGas, usedGas *big.Int, err error) {
	if err = st.preCheck(); err != nil {
		return
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	fmt.Printf("TransitionDb : StateTransition st.msg.Governance() = %b\n", st.msg.Governance())    // yhheo

	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	// Pay intrinsic gas
	// TODO convert to uint64
	intrinsicGas := IntrinsicGas(st.data, contractCreation, homestead)

	log.Info("TransitionDb","intrinsicGas", intrinsicGas) // KBW

	if intrinsicGas.BitLen() > 64 {
		return nil, nil, nil, vm.ErrOutOfGas
	}
	if err = st.useGas(intrinsicGas.Uint64()); err != nil {
		return nil, nil, nil, err
	}

	var (
		evm = st.evm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
		fee *big.Int // KBW
		org_value *big.Int // KBW
	)
	if contractCreation {
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)

		log.Info("TransitionDb","TransitionDb 1", "contractCreation","st.gas", st.gas) // KBW
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		ret, st.gas, vmerr = evm.Call(sender, st.to().Address(), st.data, st.gas, st.value)

		log.Info("TransitionDb","TransitionDb 1", "Increment the nonce","st.gas", st.gas) // KBW
	}
	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == vm.ErrInsufficientBalance {
			return nil, nil, nil, vmerr
		}
	}
	requiredGas = new(big.Int).Set(st.gasUsed())

	//fee = 0.01 Reap(10^18) + 0.01 Reap(10^18) * st.gasUsed()
	var reap_value = new(big.Int) // KBW
	reap_value.SetString("20000000000000000",10) // KBW

	var zero_value = new(big.Int) // KBW
	zero_value.SetString("0",10) // KBW

	log.Info("TransitionDb Compare","zero value", zero_value,"used value", st.gasUsed()) // KBW

	if !contractCreation && st.gasUsed().String() == zero_value.String() {
		log.Info("TransitionDb","Caculate", "Fee 1") // KBW

		fee = new(big.Int) // KBW
		fee.SetString("0",10) // KBW
	} else {
		log.Info("TransitionDb","Caculate", "Fee 2") // KBW

		fee = new(big.Int).Add(reap_value, new(big.Int).Mul(st.gasUsed(), reap_value)) // KBW
	}

	log.Info("TransitionDb","Caculate", "Fee Complte") // KBW

	org_value = new(big.Int).Mul(st.gasUsed(), st.gasPrice) // KBW

	log.Info("TransitionDb","Gas Used", st.gasUsed(), "Reap Fee Value", fee, "Ether Value", org_value, "Gas Price", st.gasPrice) // KBW

	st.refundGas()
	//st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(st.gasUsed(), st.gasPrice))
	st.state.AddBalance(st.evm.Coinbase, fee) // KBW
	//st.state.SubBalance(st.evm.Coinbase, fee)

	return ret, requiredGas, st.gasUsed(), err
}

func (st *StateTransition) refundGas() {
	// Return eth for remaining gas to the sender account,
	// exchanged at the original rate.
	sender := st.from() // err already checked
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(sender.Address(), remaining)

	// Apply refund counter, capped to half of the used gas.
	uhalf := remaining.Div(st.gasUsed(), common.Big2)
	refund := math.BigMin(uhalf, st.state.GetRefund())
	st.gas += refund.Uint64()

	var reap_value = new(big.Int) // KBW
	reap_value.SetString("20000000000000000",10) // KBW
	var fee = new(big.Int).Add(reap_value, new(big.Int).Mul(refund, reap_value)) // KBW
	var org_value = new(big.Int).Mul(st.gasUsed(), st.gasPrice) // KBW

	log.Info("refundGas","Gas Refund", refund, "Fee Value", fee, "Ether Value", org_value, "Gas Price", st.gasPrice) // KBW

	//st.state.AddBalance(sender.Address(), refund.Mul(refund, st.gasPrice))
	st.state.AddBalance(sender.Address(), fee) // KBW
	//st.state.SubBalance(sender.Address(), fee)

	//fee = 0.01 Reap(10^18) + 0.01 Reap(10^18) * cost

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(new(big.Int).SetUint64(st.gas))
}

func (st *StateTransition) gasUsed() *big.Int {
	return new(big.Int).Sub(st.initialGas, new(big.Int).SetUint64(st.gas))
}
