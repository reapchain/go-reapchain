package qManager

import (
	"bytes"
	"encoding/binary"
	//"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/ethereum/go-ethereum/consensus/quantum"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	//"math/rand"
	"os"
	"time"
)

const NodeIDBits = 512



func CheckQRNDStatus(){

	if fileExists("/Volumes/PSoC USB/" + "up.ini") {
		podc_global.QRNDDeviceStat = true
	} else {
		podc_global.QRNDDeviceStat = false
	}
}
func ConnectDB() {

	var err error
	podc_global.QManagerStorage, err = leveldb.OpenFile("level", nil)
	if err != nil{
		log.Info("DB ERROR", "err = ", err)
	}
	go StartExpirationChecker()
	CheckQRNDStatus()

	if podc_global.QRNDDeviceStat == true{

		podc_global.RandomNumbers = generateRandomNumbers()
		go StartQRNDRefresher()

	}

	go Start()

}

func IsQmanager() (isQMan bool){
	return podc_global.QManConnected
}

func StartQRNDRefresher(){

	uptimeTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-uptimeTicker.C:
			generateRandomNumbers()
		}
	}
}

func StartExpirationChecker(){

	uptimeTicker := time.NewTicker(300 * time.Second)

	for {
		select {
		case <-uptimeTicker.C:
			expirationCheck()
		}
	}
}


func expirationCheck() {
	iter :=  podc_global.QManagerStorage.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var decodedBytes podc_global.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			log.Info("Qmanager", "Decoding Error", err.Error())
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		nowtimestamp, _ := time.Parse("2006-01-02 15:04:05", timestamp )

		dbtimestamp, _ := time.Parse("2006-01-02 15:04:05", decodedBytes.Timestamp)


		diff := nowtimestamp.Sub(dbtimestamp).Seconds()
		if diff > 300 {
			log.Info("Qmanager", "Current Time", nowtimestamp)
			log.Info("Qmanager", "DB Last Updated", dbtimestamp)


			log.Info("Qmanager", "Expired Node", key)


			err = podc_global.QManagerStorage.Delete(key, nil)
			if err != nil {
				log.Info("Qmanager", "Decoding Error", err.Error())
			}
		}
	}
}


//
//func FindNode(nodeId NodeID) (found bool) {
//	//QManagerStorage, err = leveldb.OpenFile("level", nil)
//
//	//var data []byte
//	node_id_encoded,_ := rlp.EncodeToBytes(nodeId.String())
//	//log.Info("FIND ID", "bootnode = ", nodeId)
//	//log.Info("Encoded ID", "rlp = ", node_id_encoded)
//
//	found_value := qManager.Find(node_id_encoded)
//
//	return found_value
//}


func FindNode(nodeAddress string) ( found bool) {
	node_address_encoded,_ := rlp.EncodeToBytes(nodeAddress)

	//ecies.Encrypt()


	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if podc_global.QManConnected {
		//var nodeIDString string
		//decode_err := rlp.DecodeBytes(nodeId, &nodeIDString)
		//
		//if decode_err != nil{
		//	log.Info("QManager", "RLP Decode Error = ", decode_err)
		//	return
		//}
		//
		//log.Info("Decoded ID", "qman = ", nodeIDString)

		foundNode, err := podc_global.QManagerStorage.Get(node_address_encoded, nil)
		if err != nil {
			log.Info("QManager", "DB --", "Node Not Found")

			return false
		}

		var decodedBytes  podc_global.QManDBStruct
		DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
		if DecodeErr != nil {
			log.Info("Qmanager", "Decoding Error", DecodeErr.Error())
		}

		var encodedStruct *podc_global.QManDBStruct
		initBytes, err := rlp.EncodeToBytes(encodedStruct)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)

		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		encodedStruct = &podc_global.QManDBStruct{ID: decodedBytes.ID,  Address: decodedBytes.Address, Timestamp: timestamp, Tag: decodedBytes.Tag }

		//t, _ := time.Parse("2006-01-02 15:04:05", encodedStruct.Timestamp)
		//fmt.Println(t)
		initBytes, err = rlp.EncodeToBytes(encodedStruct)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)

		}

		saved :=  SaveToDB(node_address_encoded, initBytes)
		if !saved {
			log.Info("QManager DB Save", "err --", "UPDATE ERROR")
		}


		//fmt.Println(data)

		//dec := gob.NewDecoder(bytes.NewReader(data))
		//
		//dec.Decode(node)

		//defer QManagerStorage.Close()
		return true

	}

	return true

}


func  Save( dbStruct podc_global.QManDBStruct) (saved bool) {

	encodedAddress,_ := rlp.EncodeToBytes(dbStruct.Address)


	var encodedStruct *podc_global.QManDBStruct
	initBytes, err := rlp.EncodeToBytes(encodedStruct)

	if err != nil {
		log.Info("QManager DB Save", "err --", err)

	}


	encodedStruct = &podc_global.QManDBStruct{ID: dbStruct.ID,  Address: dbStruct.Address, Timestamp: dbStruct.Timestamp, Tag: dbStruct.Tag}
	//t, _ := time.Parse("2006-01-02 15:04:05", encodedStruct.Timestamp)
	//fmt.Println(t)
	initBytes, err = rlp.EncodeToBytes(encodedStruct)

	if err != nil {
		log.Info("QManager DB Save", "err --", err)

	}


	saved =  SaveToDB(encodedAddress, initBytes)
	if !saved {
		return false
	}

	return true
}


func SaveToDB(Address []byte, NodeDetails []byte) ( saved bool) {
	if podc_global.QManConnected {
		err := podc_global.QManagerStorage.Put(Address, NodeDetails, nil)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)
			return false

		}
		return true

	}
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateRandomNumbers() (RandomNumbers []uint64) {
	generatedBUFFER := quantum.GenerateQrnd()

	one := make([]byte, 8)
	var RandomNums []uint64
	counter := 0
	for i:=0; i<len(generatedBUFFER) - 4; i++{

		one[counter] = generatedBUFFER[i]
		counter = counter + 1
		if counter == 4{
			RandomNums = append(RandomNums, binary.LittleEndian.Uint64(one))
			counter = 0
		}
	}

	return RandomNums
}