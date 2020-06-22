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
	"github.com/ethereum/go-ethereum/common"
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

func generateExtraData() []ValidatorInfo{

	iter := qManager.QManagerStorage.NewIterator(nil, nil)
	var extra []ValidatorInfo  //slice의 경우 사용시 메모리를 잡아줘야 하는데 , 현재 없음.

    //make slice memory... 할것,,
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


		var num uint64

		if qManager.QRNDDeviceStat == true{
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(12280)
			num = qManager.RandomNumbers[randomIndex]


		} else {
			num = rand.Uint64()
		}

		validatorInfo := ValidatorInfo{}
		validatorInfo.Address = decodedAddress
		validatorInfo.Qrnd = num
		validatorInfo.Tag = podc.Candidate

		extra = append(extra, validatorInfo)
		//i++
	}

	return extra

}
//For Qmanager, event handler to receive msg from geth
func (c *core) handleExtraData(msg *message, src podc.Validator) error {
	if qManager.QManConnected{

		log.Info("Received EXTRA DATA REQUEST from geth")

		Counter = Counter + 1
		log.Info("Round", "Count: ", Counter)

		var extra []ValidatorInfo
		for {
			extra = generateExtraData()
			completed := false
			divisor := rand.Intn(50) + 1

			index := 0
			for index < len(extra) {

				if !c.valSet.IsProposer( extra[index].Address) && c.qmanager != extra[index].Address {
					randomNumber := extra[index].Qrnd
					if randomNumber%uint64(divisor) == 0 {
						extra[index].Tag = podc.Coordinator
						log.Info("Qmanager", "Random Coordinator Selected", extra[index].Address.String())

						completed = true
						Divisor = divisor
						break
					}
				}
				index++
			}
			if completed{
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
			Msg: extraDataJson,
		}, src.Address())

	}
	return nil
}

func (c *core) CoordinatorConfirmation(msg *message, src podc.Validator) error {
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