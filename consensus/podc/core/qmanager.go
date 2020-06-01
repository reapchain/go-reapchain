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
	"bytes"
	"encoding/binary"
	"reflect"

	//"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/consensus/quantum"
	"os"
	"encoding/json"
	"github.com/ethereum/go-ethereum/qManager"
	"math/rand"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

var (

	QManagerStorage *leveldb.DB
	Counter int
	Divisor int

)
/*
type {
	NodeID [NodeIDBits / 8]byte

	qManagerNodes struct {
		ID      NodeID
		pubKey  *ecdsa.PublicKey
		address common.Address
	}

*/
func generateExtraData() []ValidatorInfo{

	iter := qManager.QManagerStorage.NewIterator(nil, nil)
	var extra []ValidatorInfo
	//var i = 0
	//flag := false

	for iter.Next() {
		//key := iter.Key()
		value := iter.Value()
		//log.Info("KEY & Val", "key:", key, "value: ", value)

		var decodedBytes qManager.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			log.Info("Qmanager", "Decoding Error", err.Error())

		}

		decodedAddress :=  common.HexToAddress(decodedBytes.Address)
		//log.Info("Qmanager", "Address", decodedAddress)

		var num uint64

		if qManager.QRNDDeviceStat == true{
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(12280)
			num = qManager.RandomNumbers[randomIndex]
			//quant := quantum.GenerateQrnd()
			////fmt.Println(quant)
			//num = binary.LittleEndian.Uint64(quant)
			////fmt.Println(quant)
			////fmt.Println(num)
			//log.Info("Qmanager", "Quantum Number", num)

		} else {
			num = rand.Uint64()
			//log.Info("Qmanager", "Pusedo Quantum Number", num)
		}

		validatorInfo := ValidatorInfo{}
		validatorInfo.Address = decodedAddress
		validatorInfo.Qrnd = num
		validatorInfo.Tag = podc.Candidate

		//if i == 0 {
		//	if !c.valSet.IsProposer( decodedAddress) {
		//		validatorInfo.Tag = podc.Coordinator
		//	} else {
		//		flag = true
		//		validatorInfo.Tag = podc.Candidate
		//	}
		//} else if i == 1 {
		//	if flag {
		//		validatorInfo.Tag = podc.Coordinator
		//	} else {
		//		validatorInfo.Tag = podc.Candidate
		//	}
		//} else {
		//	validatorInfo.Tag = podc.Candidate
		//}
		extra = append(extra, validatorInfo)
		//i++
	}

	return extra

}
func (c *core) handleExtraData(msg *message, src podc.Validator) error {
	if (reflect.DeepEqual(c.qmanager, c.Address())) { //if I'm Qmanager
		if qManager.QManConnected {

			log.Info("EXTRA DATA REQUEST")
			//log.Info("Requesting Source", "from", src)
			Counter = Counter + 1
			log.Info("Round", "Count: ", Counter)

			var extra []ValidatorInfo
			for {
				extra = generateExtraData()
				completed := false
				divisor := rand.Intn(50) + 1

				index := 0
				for index < len(extra) {

					if !c.valSet.IsProposer(extra[index].Address) && c.qmanager != extra[index].Address {
						randomNumber := extra[index].Qrnd
						if randomNumber%uint64(divisor) == 0 {
							//log.Info("COORDINATOR", "Random Divisor", divisor)
							//log.Info("COORDINATOR", "Random Number", randomNumber)
							extra[index].Tag = podc.Coordinator
							log.Info("Qmanager", "Random Coordinator Selected", extra[index].Address.String())
							index = len(extra)
							completed = true
							Divisor = divisor
						}
					}
					index++
				}
				if completed {
					break
				}
			}

			log.Info("ExtraData list", "extradata", extra)

			//defer db.Close()
			extraDataJson, err := json.Marshal(extra)
			if err != nil {
				log.Error("Failed to encode JSON", err)
			}

			c.send(&message{
				Code: msgExtraDataSend,
				Msg:  extraDataJson,
			}, src.Address())
			// Decode commit message
			//fmt.Println("EXTRA DATA HANDLE")
			//fmt.Println(src)

		}
	}
	return nil
}

func (c *core) CoordinatorConfirmation(msg *message, src podc.Validator) error {
	////logger := c.logger.New("EXTRA DATA")
	//log.Trace("EXTRA DATA SENT DATA")
	////logger := c.logger.New()
	////logger.Info("EXTRA DATA REQUEST")
	////log.Debug("from", src)
	//log.Debug("Requesting Source", "from", src)
	//log.Debug("ExtraDataMessage", "from", msg)

	CoordiQRND := binary.LittleEndian.Uint64(msg.Msg)

	if CoordiQRND%uint64(Divisor) == 0 {
		log.Info("QManager", "Coordinator Confirm Status: ", true)

		c.send(&message{
			Code: msgCoordinatorConfirmSend,
			Msg: []byte("Coordinator Confirmed"),
		}, src.Address())

	} else{
		log.Info("QManager Error", "Coordinator Confirm Status: ", false)
	}



	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}