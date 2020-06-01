package qManager

import (
	"bytes"
	"encoding/binary"
	//"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	//"math/rand"
	"net"
	"os"
	"time"
	"github.com/ethereum/go-ethereum/consensus/quantum"

)

const NodeIDBits = 512


var (

	QManagerStorage *leveldb.DB
	QManConnected bool
	QManagerAddress *net.UDPAddr
	BootNodeReady bool
	QRNDDeviceStat bool
 	RandomNumbers []uint64
	BootNodePort int
	IsBootNode bool
)

type (
	QManDBStruct struct {
		ID      string
		Address  string
		Timestamp	string
	}
)

func CheckQRNDStatus(){

	if fileExists("/Volumes/PSoC USB/" + "up.ini") {
		QRNDDeviceStat = true
	} else {
		QRNDDeviceStat = false
	}
}
func ConnectDB() {

	var err error
	QManagerStorage, err = leveldb.OpenFile("level", nil)
	if err != nil{
		log.Info("DB ERROR", "err = ", err)
	}
	go StartExpirationChecker()
	CheckQRNDStatus()

	if QRNDDeviceStat == true{

		RandomNumbers = generateRandomNumbers()
		go StartQRNDRefresher()

	}

}

func IsQmanager() (isQMan bool){
	return QManConnected
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
	iter :=  QManagerStorage.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var decodedBytes QManDBStruct
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


			err = QManagerStorage.Delete(key, nil)
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


func FindNode(nodeId string) ( found bool) {
	node_id_encoded,_ := rlp.EncodeToBytes(nodeId)


	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if QManConnected {
		//var nodeIDString string
		//decode_err := rlp.DecodeBytes(nodeId, &nodeIDString)
		//
		//if decode_err != nil{
		//	log.Info("QManager", "RLP Decode Error = ", decode_err)
		//	return
		//}
		//
		//log.Info("Decoded ID", "qman = ", nodeIDString)

		foundNode, err := QManagerStorage.Get(node_id_encoded, nil)
		if err != nil {
			log.Info("QManager", "DB --", "Node Not Found")

			return false
		}

		var decodedBytes  QManDBStruct
		DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
		if DecodeErr != nil {
			log.Info("Qmanager", "Decoding Error", DecodeErr.Error())
		}

		var encodedStruct *QManDBStruct
		initBytes, err := rlp.EncodeToBytes(encodedStruct)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)

		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		encodedStruct = &QManDBStruct{ID: decodedBytes.ID,  Address: decodedBytes.Address, Timestamp: timestamp }

		//t, _ := time.Parse("2006-01-02 15:04:05", encodedStruct.Timestamp)
		//fmt.Println(t)
		initBytes, err = rlp.EncodeToBytes(encodedStruct)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)

		}

		saved :=  SaveToDB(node_id_encoded, initBytes)
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


func ( n QManDBStruct) Save() (saved bool) {

	encodedID,_ := rlp.EncodeToBytes(n.ID)


	var encodedStruct *QManDBStruct
	initBytes, err := rlp.EncodeToBytes(encodedStruct)

	if err != nil {
		log.Info("QManager DB Save", "err --", err)

	}


	encodedStruct = &QManDBStruct{ID: n.ID,  Address: n.Address, Timestamp: n.Timestamp }
	//t, _ := time.Parse("2006-01-02 15:04:05", encodedStruct.Timestamp)
	//fmt.Println(t)
	initBytes, err = rlp.EncodeToBytes(encodedStruct)

	if err != nil {
		log.Info("QManager DB Save", "err --", err)

	}


	saved =  SaveToDB(encodedID, initBytes)
	if !saved {
		return false
	}

	return true
}


func SaveToDB(ID []byte, NodeDetails []byte) ( saved bool) {
	if QManConnected {
		err := QManagerStorage.Put(ID, NodeDetails, nil)

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