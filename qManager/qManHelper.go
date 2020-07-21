package qManager

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/quantum"
	"runtime"

	//"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	//"math/rand"
	"os"
	"time"
)

const NodeIDBits = 512



func CheckQRNDStatus(){

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		fmt.Println("Windows")
	case "darwin":
		podc_global.QRNGFilePrefix = "/Volumes/PSoC USB/"
	case "linux":
		user := os.Getenv("USERNAME")
		podc_global.QRNGFilePrefix = "/media/"+ user +"/E8EE-1C60/"
	default:
		user := os.Getenv("USERNAME")
		podc_global.QRNGFilePrefix = "/media/"+ user +"/E8EE-1C60/"
	}

	log.Info("Qmanager", "QRND = ", podc_global.QRNGFilePrefix)
	if fileExists(podc_global.QRNGFilePrefix + "up.ini") {
		podc_global.QRNDDeviceStat = true
		//log.Info("QRND", "Buffer: ", "GENERATING NUMS")
		podc_global.RandomNumbers = generateRandomNumbers()
		go StartQRNDRefresher()

	} else {
		podc_global.QRNDDeviceStat = false
	}
}

func ConnectDB(){
	var err error
	podc_global.QManagerStorage, err = leveldb.OpenFile("level", nil)
	if err != nil{
		log.Info("DB ERROR", "err = ", err)
	}
}
func CloseDB(){
	podc_global.QManagerStorage.Close()
}

func QmanInit() {


	go StartExpirationChecker()
	CheckQRNDStatus()


	//if podc_global.QRNDDeviceStat == true{
	//
	//	//podc_global.RandomNumbers = generateRandomNumbers()


	//}

	//go Start()

}

func IsQmanager() (isQMan bool){
	return podc_global.QManConnected
}

func StartQRNDRefresher(){


	uptimeTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-uptimeTicker.C:
			podc_global.RandomNumbers = generateRandomNumbers()
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

	ConnectDB()
	log.Info("Qmanager", "DB Status", "1. Connected")
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


			log.Info("Qmanager", "Expired Node", string(key))


			err = podc_global.QManagerStorage.Delete(key, nil)
			if err != nil {
				log.Info("Qmanager", "Decoding Error", err.Error())
			}
		}
	}
	CloseDB()
	log.Info("Qmanager", "DB Status", "1. Disconnected")

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


func UpdateSenatorCandidateNodes() {
	if podc_global.QManConnected {
		ConnectDB()
		log.Info("Qmanager", "DB Status", "2. Connected")


		for _, element := range podc_global.GovernanceList {
			nodeAddress := element.Validator

			node_address_encoded,_ := rlp.EncodeToBytes(nodeAddress)
			foundNode, err := podc_global.QManagerStorage.Get(node_address_encoded, nil)
			if err != nil {
				log.Info("QManager", "DB --", "Node Not Found")

			}

			if foundNode != nil{
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
				tagConverted, _ := math.ParseUint64(element.Tag)
				encodedStruct = &podc_global.QManDBStruct{ID: decodedBytes.ID,  Address: decodedBytes.Address,
					Timestamp: decodedBytes.Timestamp, Tag:tagConverted }

				initBytes, err = rlp.EncodeToBytes(encodedStruct)

				if err != nil {
					log.Info("QManager", "RLP Error --", err)

				}
				saveError := podc_global.QManagerStorage.Put(node_address_encoded, initBytes, nil)

				if saveError != nil {
					log.Info("QManager DB Save", "err --", saveError)
				}

				//
				//saved :=  SaveToDB(node_address_encoded, initBytes)
				//if !saved {
				//	log.Info("QManager DB Save", "err --", "UPDATE ERROR")
				//}
			}
		}
		CloseDB()
		log.Info("Qmanager", "DB Status", "2. Disconnected")

	}
}


func FindNode(nodeAddress string) ( found bool) {
	node_address_encoded,_ := rlp.EncodeToBytes(nodeAddress)

	//ecies.Encrypt()


	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if podc_global.QManConnected {
		ConnectDB()
		log.Info("Qmanager", "DB Status", "3. Connected")

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
			CloseDB()
			log.Info("Qmanager", "DB Status", "3. Disconnected")

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




		updateError := podc_global.QManagerStorage.Put(node_address_encoded, initBytes, nil)

		if updateError != nil {
			log.Info("QManager DB Save", "UPDATE ERROR --", updateError)
		}

		//saved :=  SaveToDB(node_address_encoded, initBytes)
		//if !saved {
		//	log.Info("QManager DB Save", "err --", "UPDATE ERROR")
		//}


		//fmt.Println(data)

		//dec := gob.NewDecoder(bytes.NewReader(data))
		//
		//dec.Decode(node)

		//defer QManagerStorage.Close()
		CloseDB()
		log.Info("Qmanager", "DB Status", "3. Disconnected")

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
		ConnectDB()
		log.Info("Qmanager", "DB Status", "4. Connected")

		saveErr := podc_global.QManagerStorage.Put(Address, NodeDetails, nil)

		if saveErr != nil {
			log.Info("QManager DB Save", "err --", saveErr)
			CloseDB()
			log.Info("Qmanager", "DB Status", "4. Disconnected")

			return false

		}

		CloseDB()
		log.Info("Qmanager", "DB Status", "4. Disconnected")

		return true

	}
	return false
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
	//log.Info("QRND", "Buffer: ", generatedBUFFER)


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

func GetIterator() () {


}
