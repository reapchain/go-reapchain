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
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/quantum"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	//"math/rand"
)
//const NodeIDBits = 512

var (

	QManagerStorage *leveldb.DB
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

func (c *core) handleExtraData(msg *message, src istanbul.Validator) error {
	if common.QManConnected{


	//logger := c.logger.New("EXTRA DATA")
	log.Info("EXTRA DATA REQUEST")


	//logger := c.logger.New()
	//logger.Info("EXTRA DATA REQUEST")
	//log.Debug("from", src)
	log.Info("Requesting Source", "from", src)


	//QManagerStorage, _ = leveldb.OpenFile("level", nil)
	//db, err := leveldb.OpenFile("level", nil)
	//if err != nil{
	//	log.Info("DB ERROR", "err = ", err)
	//}

	//QManagerStorage, _ = leveldb.OpenFile("level", nil)



	iter := common.QManagerStorage.NewIterator(nil, nil)
	var extra []ValidatorInfo
	var i = 0
	flag := false

	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		log.Info("KEY & Val", "key:", key, "value: ", value)
		//var qNode qManagerNodes
		//errs := json.Unmarshal(iter.Value(), &qNode)
		//if errs != nil {
		//	fmt.Println("error:", err)
		//}
		//log.Debug("Data", "address:", qNode.address, "ID: ", qNode.ID)
		var decodedBytes common.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		} else {
			fmt.Printf("Decoded value: %#v\n", decodedBytes)
		}

		decodedAddress := common.HexToAddress(decodedBytes.Address)
		//decodedNodeID,_ := discover.HexID(decodedBytes.ID)

		quant := quantum.GenerateQrnd()
		//fmt.Println(quant)
		num := binary.LittleEndian.Uint64(quant)
		fmt.Println(num)


		validatorInfo := ValidatorInfo{}
		validatorInfo.Address = decodedAddress
		validatorInfo.Qrnd = num


		//validatorInfo := ValidatorInfo{}
		//validatorInfo.Address = common.HexToAddress(string(value))
		//validatorInfo.Qrnd = rand.Uint64()

		if i == 0 {
			if !c.valSet.IsProposer(decodedAddress) {
				validatorInfo.Tag = istanbul.Coordinator
			} else {
				flag = true
				validatorInfo.Tag = istanbul.Candidate
			}
		} else if i == 1 {
			if flag {
				validatorInfo.Tag = istanbul.Coordinator
			} else {
				validatorInfo.Tag = istanbul.Candidate
			}
		} else {
			validatorInfo.Tag = istanbul.Candidate
		}
		extra = append(extra, validatorInfo)
		i++
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
	// Decode commit message
	//fmt.Println("EXTRA DATA HANDLE")
	//fmt.Println(src)

	}
	return nil
}

func (c *core) handleSentData(msg *message, src istanbul.Validator) error {
	//logger := c.logger.New("EXTRA DATA")
	log.Trace("EXTRA DATA SENT DATA")


	//logger := c.logger.New()
	//logger.Info("EXTRA DATA REQUEST")
	//log.Debug("from", src)
	log.Debug("Requesting Source", "from", src)
	log.Debug("ExtraDataMessage", "from", msg)





	return nil
}

