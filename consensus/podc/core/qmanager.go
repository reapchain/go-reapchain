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
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/podc"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"github.com/ethereum/go-ethereum/rlp"
	"math/rand"
	"os"
	"time"
)

var (

	Counter int
	Divisor int

)

func generateExtraData(dbData []ValidatorInfo) []ValidatorInfo{

	//qManager.ConnectDB()
	//log.Info("Qmanager", "DB Status", "4. Connected")

	var extra []ValidatorInfo
	for index, validator := range dbData {

		var num uint64

		if podc_global.QRNDDeviceStat == true{
			log.Info("QRND " + string(index), " Random Nums" , podc_global.QRNDDeviceStat)
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(12280)
			num = podc_global.RandomNumbers[randomIndex]


		} else {
			log.Info("Suedo Random "  + string(index) , " Random Nums", podc_global.QRNDDeviceStat)
			num = rand.Uint64()
		}

		validator.Qrnd = num
		extra = append(extra, validator)
	}

	//iter := snapshot.NewIterator(nil, nil)


    //make slice memory... 할것,,
	//for iter.Next() {
	//	//key := iter.Key()
	//	value := iter.Value()
	//	//log.Info("KEY & Val", "key:", key, "value: ", value)
	//
	//	var decodedBytes podc_global.QManDBStruct
	//	err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
	//	if err != nil {
	//		log.Info("Qmanager", "Decoding Error", err.Error())
	//
	//	}
	//
	//	decodedAddress :=  common.HexToAddress(decodedBytes.Address)
	//
	//
	//	var num uint64
	//
	//	if podc_global.QRNDDeviceStat == true{
	//		log.Info("QRND ", "Random Nums", podc_global.QRNDDeviceStat)
	//		rand.Seed(time.Now().UnixNano())
	//		randomIndex := rand.Intn(12280)
	//		num = podc_global.RandomNumbers[randomIndex]
	//
	//
	//	} else {
	//		log.Info("Suedo Random ", "Random Nums", podc_global.QRNDDeviceStat)
	//		num = rand.Uint64()
	//	}
	//
	//	validatorInfo := ValidatorInfo{}
	//	validatorInfo.Address = decodedAddress
	//	validatorInfo.Qrnd = num
	//	validatorInfo.Tag = podc.Tag(decodedBytes.Tag)
	//
	//	extra = append(extra, validatorInfo)
	//	//i++
	//}
    //log.Info("generateExtraData:", "extradata_size", len(extra))
    //qManager.CloseDB()
	//log.Info("Qmanager", "DB Status", "4. Disconnected")

	return extra

}
//For Qmanager, event handler to receive msg from geth
func (c *core) handleExtraData(msg *message, src podc.Validator) error {
	if podc_global.QManConnected{

		log.Info("Received EXTRA DATA REQUEST from geth")

		Counter = Counter + 1
		log.Info("Round", "Count: ", Counter)

		var extra []ValidatorInfo
		qManager.ConnectDB()
		//log.Info("Qmanager", "DB Status", "5. Connected")
		//
		//dbSnapShot, _ := podc_global.QManagerStorage.GetSnapshot()
		//
		//qManager.CloseDB()
		//log.Info("Qmanager", "DB Status", "5. Disconnected")

		var dbData []ValidatorInfo  //slice의 경우 사용시 메모리를 잡아줘야 하는데 , 현재 없음.
		iter := podc_global.QManagerStorage.NewIterator(nil, nil)

		for iter.Next() {
			value := iter.Value()
			var decodedBytes podc_global.QManDBStruct
			err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
			if err != nil {
				log.Info("Qmanager", "Decoding Error", err.Error())
			}
			decodedAddress :=  common.HexToAddress(decodedBytes.Address)
			validatorInfo := ValidatorInfo{}
			validatorInfo.Address = decodedAddress
			validatorInfo.Qrnd = 1
			validatorInfo.Tag = podc.Tag(decodedBytes.Tag)
			dbData = append(dbData, validatorInfo)
		}

		qManager.CloseDB()
		log.Info("Qmanager", "DB Status", "5. Disconnected")

		for {
			log.Info("Qmanager", "Generating Random Numbers", "Outerloop")
			extra = generateExtraData(dbData)
			completed := false
			divisor := rand.Intn(50) + 1

			index := 0
			for index < len(extra) {
				log.Info("Qmanager", "Generating Random Numbers", "InnerLoop")


				if !c.valSet.IsProposer( extra[index].Address) && c.qmanager != extra[index].Address {
					randomNumber := extra[index].Qrnd
					if randomNumber%uint64(divisor) == 0 {
						extra[index].Tag = podc.Coordinator
						log.Info("Qmanager", "Random Coordinator Selected", extra[index].Address.String())

						index = len(extra)
						completed = true
						Divisor = divisor
					}
				}
				//log.Info("ExtraData list", "Address", extra[index].Address , "Qrnd", extra[index].Qrnd, "Tag",  extra[index].Tag)
				index++
			}
			if completed{
				break
			}
			if len(extra) == 0 {
				break
			}
		}


		//log.Info("ExtraData list", "extradata", extra)

		//defer db.Close()
		log.Info("QManager", "ExtraData Length: ", len(extra))
		log.Info("QManager", "ExtraData: ", extra)
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