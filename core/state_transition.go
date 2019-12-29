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

	Governance() bool 	// yhheo
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
		igas.SetUint64(params.TxGas)
	}
	//fmt.Printf("IntrinsicGas : base gas = %d\n", igas)	// yhheo

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

	mgval := st.calcFee(mgas)   // yhheo - new(big.Int).Mul(mgas, st.gasPrice)

	var (
		state  = st.state
		sender = st.from()
	)
	if !st.msg.Governance() &&!st.coinTransfer() { // yhheo
		if state.GetBalance(sender.Address()).Cmp(mgval) < 0 {
			return errInsufficientBalanceForGas
		}
	}
	if err := st.gp.SubGas(mgas); err != nil {
		return err
	}
	st.gas += mgas.Uint64()

	st.initialGas.Set(mgas)
    // yhheo - begin
    if st.msg.Governance() {
        return nil
    }
    if st.coinTransfer() {
        return nil
    }
    // yhheo - end
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
	fmt.Println("TransitionDb : st.msg.Governance() =", st.msg.Governance())	// yhheo
	if err = st.preCheck(); err != nil {
		return
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	fmt.Printf("TransitionDb : st.evm.Coinbase = %x\n", st.evm.Coinbase)	// yhheo
	//fmt.Printf("TransitionDb : st.gas = %d\n st.gasPrice = %d\n", st.gas, st.gasPrice)	// yhheo

	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	//fmt.Printf("TransitionDb : homestead = %t\n contractCreation = %t\n", homestead, contractCreation)         // yhheo

	// Pay intrinsic gas
	// TODO convert to uint64
	intrinsicGas := IntrinsicGas(st.data, contractCreation, homestead)
	if intrinsicGas.BitLen() > 64 {
		return nil, nil, nil, vm.ErrOutOfGas
	}
	if err = st.useGas(intrinsicGas.Uint64()); err != nil {
		return nil, nil, nil, err
	}

    //fmt.Printf("TransitionDb : intrinsicGas = %d\n", intrinsicGas)        // yhheo

	var (
		evm = st.evm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
        fee *big.Int    		// yhheo
	)
	if contractCreation {
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
		//fmt.Printf("TransitionDb : evm.Create (st.gas) = %d\n", st.gas)	// yhheo
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		ret, st.gas, vmerr = evm.Call(sender, st.to().Address(), st.data, st.gas, st.value)
		//fmt.Printf("TransitionDb : evm.Call (st.gas) = %d\n", st.gas)	// yhheo
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
    // yhheo - begin
	//fmt.Printf("TransitionDb : st.coinTransfer() = %t\n", st.coinTransfer())  // yhheo
    if st.msg.Governance() || st.coinTransfer() {
		requiredGas = new(big.Int).SetUint64(0)
		st.gas = st.initialGas.Uint64()
	} else {
		requiredGas = new(big.Int).Set(st.gasUsed())
	}
    // yhheo - end

	st.refundGas()

    // yhheo - begin
	requiredGas = new(big.Int).Set(st.gasUsed())
	fee = st.calcFee(st.gasUsed())
    //fmt.Printf("TransitionDb : fee = %d\n", fee)

	st.state.AddBalance(st.evm.Coinbase, fee)   // new(big.Int).Mul(st.gasUsed(), st.gasPrice)
    // yhheo - end

    //fmt.Printf("TransitionDb : requiredGas = %d\n st.gasUsed() = %d\n", requiredGas, st.gasUsed())
	return ret, requiredGas, st.gasUsed(), err
}

func (st *StateTransition) refundGas() {
	// Return eth for remaining gas to the sender account,
	// exchanged at the original rate.
	sender := st.from() // err already checked
    // yhheo - begin
	var remaining	*big.Int
    if !st.msg.Governance() && !st.coinTransfer() {
	//	remaining = new(big.Int).SetUint64(0)
	//	fmt.Printf("refundGas : remaining = %d\n st.gas = %d\n", remaining, st.gas)
	//} else {
		remaining = st.calcFee(new(big.Int).SetUint64(st.gas)) // new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
		st.state.AddBalance(sender.Address(), remaining)
		//fmt.Printf("refundGas : remaining = %d\n st.gas = %d\n", remaining, st.gas)
		// yhheo - end

		// Apply refund counter, capped to half of the used gas.
		uhalf := remaining.Div(st.gasUsed(), common.Big2)
		refund := math.BigMin(uhalf, st.state.GetRefund())
		st.gas += refund.Uint64()
		//fmt.Printf("refundGas : st.state.GetRefund() = %d\n refund gas = %d\n", st.state.GetRefund(), refund.Uint64())	// yhheo

		// yhheo - begin
		refundFee := st.calcFee(refund)   // refund.Mul(refund, st.gasPrice)
		//fmt.Printf("refundGas : refundFee = %d\n st.gas = %d\n", refundFee, st.gas)
		st.state.AddBalance(sender.Address(), refundFee)
    }
	// yhheo - end

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(new(big.Int).SetUint64(st.gas))
}

func (st *StateTransition) gasUsed() *big.Int {
	return new(big.Int).Sub(st.initialGas, new(big.Int).SetUint64(st.gas))
}

// yhheo - begin
func (st *StateTransition) calcFee(gas *big.Int) *big.Int {
    var (
        zeroGas = new(big.Int)
        zeroFee = new(big.Int)
        baseFee = new(big.Int)
        opcdFee = new(big.Int)
    )
    if st.msg.Governance() {
        return zeroFee.SetUint64(0)
    }
    //if gas == zeroGas.SetUint64(0)  {
    if gas.Cmp(zeroGas.SetUint64(0)) == 0 {
        return zeroFee.SetUint64(0)
    }
    baseFee.SetString("20000000000000000",10)   // 0.02  reap = 10^16 wei
    opcdFee.SetString( "3000000000000000",10)   // 0.003 reap = 10^15 wei

    // reapchain fee = 0.02 Reap(10^16) + ((opcode연산:gas) * 0.003 Reap(10^15))
    return new(big.Int).Add(baseFee, new(big.Int).Mul(gas, opcdFee))
}

func (st *StateTransition) coinTransfer() bool {
    return len(st.msg.Data()) == 0
}
// yhheo - end
