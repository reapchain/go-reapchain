package qmanager

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"github.com/ethereum/go-ethereum/qmanager/quantum"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"runtime"
	"time"
)

var IsInUse = false

func CheckQRNGStatus(){
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		fmt.Println("Windows")
	case "darwin":
		global.QRNGFilePrefix = "/Volumes/PSoC USB/"
	case "linux":
		user := os.Getenv("USER")
		global.QRNGFilePrefix = "/media/"+ user +"/E8EE-1C60/"
	default:
		user := os.Getenv("USER")
		global.QRNGFilePrefix = "/media/"+ user +"/E8EE-1C60/"
	}

	log.Info("QRNG Device ", "Path = ", global.QRNGFilePrefix)
	if fileExists(global.QRNGFilePrefix + "up.ini") {
		global.QRNGDeviceStat = true
		//log.Info("QRND", "Buffer: ", "GENERATING NUMS")
		var expectedErr error
		global.RandomNumbers, expectedErr = generateRandomNumbers()
		if expectedErr == nil{
			go StartQRNGRefresher()
		}else{
			global.QRNGDeviceStat = false
		}
	} else {
		global.QRNGDeviceStat = false
	}
}

func ConnectDB(){
	//log.Info("Qmanager", "IsInUse OPEN = ", IsInUse)
	if IsInUse == false {
		IsInUse = true
		//log.Info("Qmanager", "DB = ", "STARTING---------------------------------------")
		var err error
		global.QManagerStorage, err = leveldb.OpenFile("qman-" + DBName, nil)
		if err != nil{
			log.Info("DB ERROR", "err = ", err)
			IsInUse = false
		}
	}else{
		//log.Info("Qmanager", "DB = ", "WAITING---------------------------------------")
		time.Sleep(100 *time.Millisecond)
		ConnectDB()
	}
}
func CloseDB(){
	//log.Info("Qmanager", "IsInUse CLOSE = ", IsInUse)
	global.QManagerStorage.Close()
	IsInUse = false
}

func InitializeQManager() {
	go StartExpirationChecker()
	CheckQRNGStatus()
}

func IsQmanager() (isQMan bool){
	return global.QManConnected
}

func StartQRNGRefresher(){
	uptimeTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			var expectedErr error
			global.RandomNumbers, expectedErr = generateRandomNumbers()
			if expectedErr != nil{
				global.QRNGDeviceStat = false
			}
		}
	}
}

func StartExpirationChecker(){

	uptimeTicker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			expirationCheck()
		}
	}
}

func GetDBData(){
	var tempDBDataList []global.QManDBStruct
	iter :=  global.QManagerStorage.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		var decodedBytes global.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			log.Info("Qmanager", "Decoding Error", err.Error())
		}
		tempDBDataList = append(tempDBDataList, decodedBytes)
	}

	global.DBDataList = tempDBDataList
}

func expirationCheck() {

	//log.Info("Qmanager", "DB Status", "1. Connected")
	ConnectDB()

	iter :=  global.QManagerStorage.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var decodedBytes global.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			log.Info("Qmanager RLP", "Decoding Error", err.Error())
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		nowtimestamp, _ := time.Parse("2006-01-02 15:04:05", timestamp )
		dbtimestamp, _ := time.Parse("2006-01-02 15:04:05", decodedBytes.Timestamp)

		diff := nowtimestamp.Sub(dbtimestamp).Seconds()
		if diff > 30 {
			log.Info("Qman Expired Node", "Current Time", nowtimestamp)
			log.Info("Qman Expired Node", "DB Last Updated", dbtimestamp)


			log.Info("Qman Expired Node", "Expired Node", string(key))


			err = global.QManagerStorage.Delete(key, nil)
			GetDBData()
			if err != nil {
				log.Info("Qmanager RLP", "Decoding Error", err.Error())
			}
		}
	}
	CloseDB()
	//log.Info("Qmanager", "DB Status", "1. Disconnected")

}

func UpdateSenatorCandidateNodes() {
	if global.QManConnected {
		//log.Info("Qmanager", "DB Status", "2. Connected")
		ConnectDB()
		for _, element := range global.GovernanceList {
			nodeAddress := element.Validator

			node_address_encoded,_ := rlp.EncodeToBytes(nodeAddress)
			foundNode, err := global.QManagerStorage.Get(node_address_encoded, nil)
			if err != nil {
				log.Info("QManager ", "DB --", "Node Not Found")

			}

			if foundNode != nil{
				var decodedBytes global.QManDBStruct
				DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
				if DecodeErr != nil {
					log.Info("Qmanager", "Decoding Error", DecodeErr.Error())
				}

				var encodedStruct *global.QManDBStruct
				initBytes, err := rlp.EncodeToBytes(encodedStruct)

				if err != nil {
					log.Info("QManager DB Save", "err --", err)

				}
				tagConverted, _ := math.ParseUint64(element.Tag)
				encodedStruct = &global.QManDBStruct{ID: decodedBytes.ID,  Address: decodedBytes.Address,
					Timestamp: decodedBytes.Timestamp, Tag:tagConverted }

				initBytes, err = rlp.EncodeToBytes(encodedStruct)

				if err != nil {
					log.Info("QManager", "RLP Error --", err)

				}
				saveError := global.QManagerStorage.Put(node_address_encoded, initBytes, nil)

				if saveError != nil {
					log.Info("QManager DB Save", "err --", saveError)
				}
				GetDBData()

			}
		}
		CloseDB()
		//log.Info("Qmanager", "DB Status", "2. Disconnected")

	}
}


func FindNode(nodeAddress string) ( found bool) {
	node_address_encoded,_ := rlp.EncodeToBytes(nodeAddress)
	//ecies.Encrypt()
	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if global.QManConnected {
		//log.Info("Qmanager", "DB Status", "3. Connected")
		ConnectDB()
		foundNode, err := global.QManagerStorage.Get(node_address_encoded, nil)
		if err != nil {
			CloseDB()
			//log.Info("Qmanager", "DB Status", "3. Disconnected")
			log.Info("QManager ", "DB --", "Node Not Found")
			return false
		}

		var decodedBytes global.QManDBStruct
		DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
		if DecodeErr != nil {
			log.Info("Qmanager ", "Decoding Error", DecodeErr.Error())
		}

		var encodedStruct *global.QManDBStruct
		initBytes, err := rlp.EncodeToBytes(encodedStruct)

		if err != nil {
			log.Info("QManager DB Save", "err --", err)
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		encodedStruct = &global.QManDBStruct{ID: decodedBytes.ID,  Address: decodedBytes.Address, Timestamp: timestamp, Tag: decodedBytes.Tag }
		initBytes, err = rlp.EncodeToBytes(encodedStruct)
		if err != nil {
			log.Info("QManager DB Save", "err --", err)
		}
		updateError := global.QManagerStorage.Put(node_address_encoded, initBytes, nil)
		if updateError != nil {
			log.Info("QManager DB Save", "UPDATE ERROR --", updateError)
		}
		GetDBData()
		CloseDB()
		//log.Info("Qmanager", "DB Status", "3. Disconnected")

		return true
	}
	return true
}


func  Save( dbStruct global.QManDBStruct) (saved bool) {

	encodedAddress,_ := rlp.EncodeToBytes(dbStruct.Address)
	var encodedStruct *global.QManDBStruct
	initBytes, err := rlp.EncodeToBytes(encodedStruct)
	if err != nil {
		log.Info("QManager DB Save", "err --", err)
	}
	encodedStruct = &global.QManDBStruct{ID: dbStruct.ID,  Address: dbStruct.Address, Timestamp: dbStruct.Timestamp, Tag: dbStruct.Tag}
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
	if global.QManConnected {
		//log.Info("Qmanager", "DB Status", "4. Connected")
		ConnectDB()
		saveErr := global.QManagerStorage.Put(Address, NodeDetails, nil)
		if saveErr != nil {
			log.Info("QManager DB Save", "err --", saveErr)
			CloseDB()
			log.Info("Qmanager", "DB Status", "4. Disconnected")
			return false
		}
		GetDBData()
		CloseDB()
		//log.Info("Qmanager", "DB Status", "4. Disconnected")
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

func generateRandomNumbers() (RandomNumbers []uint64, err error) {
	generatedBUFFER, err := quantum.GenerateQRNGData()
	if err == nil{
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
		return RandomNums, nil
	}else{
		return nil, err
	}
}
