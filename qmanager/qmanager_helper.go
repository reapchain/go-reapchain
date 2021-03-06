package qmanager

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"github.com/ethereum/go-ethereum/qmanager/quantum"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

var IsInUse = false

func CheckQRNGStatus() {
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		fmt.Println("Windows")
	case "darwin":
		global.QRNGFilePrefix = "/Volumes/PSoC USB/"
	case "linux":
		user := os.Getenv("USER")
		global.QRNGFilePrefix = "/media/" + user + "/E8EE-1C60/"
	default:
		user := os.Getenv("USER")
		global.QRNGFilePrefix = "/media/" + user + "/E8EE-1C60/"
	}

	log.Info("QRNG Device ", "Path = ", global.QRNGFilePrefix)
	if fileExists(global.QRNGFilePrefix + "up.ini") {
		global.QRNGDeviceStat = true
		//log.Info("QRND", "Buffer: ", "GENERATING NUMS")
		var expectedErr error
		global.RandomNumbers, expectedErr = generateRandomNumbers()
		if expectedErr == nil {
			go StartQRNGRefresher()
		} else {
			global.QRNGDeviceStat = false
		}
	} else {
		global.QRNGDeviceStat = false
	}
}

func ConnectDB() {
	//log.Info("Qmanager", "IsInUse OPEN = ", IsInUse)
	if IsInUse == false {
		IsInUse = true
		//log.Info("Qmanager", "DB = ", "STARTING---------------------------------------")
		//var err error
		storage, err := leveldb.OpenFile("qman-"+DBName, nil)
		if err != nil {
			log.Info("DB ERROR", "err = ", err)
			IsInUse = false
			return
		}
		global.QManagerStorage = storage
	} else {
		//log.Info("Qmanager", "DB = ", "WAITING---------------------------------------")
		time.Sleep(100 * time.Millisecond)
		ConnectDB()
	}
}
func CloseDB() {
	//log.Info("Qmanager", "IsInUse CLOSE = ", IsInUse)
	global.QManagerStorage.Close()
	IsInUse = false
}

func InitializeQManager() {
	go StartExpirationChecker()
	CheckQRNGStatus()
	go InitialValidatorConfigParsing()
	go PeriodicValidatorConfigParsing()
}

func InitialValidatorConfigParsing() {
	CheckConfiValidators()
}

func CheckConfiValidators() {
	var ConfigSenatorList = config.Config.Senatornodes
	var ConfigCandidateList = config.Config.Candidatenodes

	if len(ConfigSenatorList) == 0 || len(ConfigCandidateList) == 0 {
		log.Error("Config.json Error", "Senator & Candidate List", "Insert Senator & Candidate into Config.Json")
	} else {
		log.Info("Parsing Config.json - Senator & Candidate List")
		var govStruct []global.GovStruct
		for _, item := range ConfigSenatorList {
			senate := global.GovStruct{Validator: item, Tag: common.Senator}
			govStruct = append(govStruct, senate)
			log.Info("Parsing Config.json", "Senator Struct", senate)
		}
		for _, item := range ConfigCandidateList {
			candidate := global.GovStruct{Validator: item, Tag: common.Candidate}
			govStruct = append(govStruct, candidate)
			log.Info("Parsing Config.json", "Candidate Struct", candidate)

		}
		global.GovernanceList = govStruct

	}
}

func StartQRNGRefresher() {
	uptimeTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			var expectedErr error
			global.RandomNumbers, expectedErr = generateRandomNumbers()
			if expectedErr != nil {
				global.QRNGDeviceStat = false
			}
		}
	}
}

func PeriodicValidatorConfigParsing() {
	uptimeTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			config.Config.GetConfig("REAPCHAIN_ENV", "SETUP_INFO")
			CheckConfiValidators()
			go UpdateSenatorCandidateNodes()
		}
	}
}

func StartExpirationChecker() {

	uptimeTicker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			expirationCheck()
		}
	}
}

func GetDBData() {
	var tempDBDataList []global.QManDBStruct
	iter := global.QManagerStorage.NewIterator(nil, nil)
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

	iter := global.QManagerStorage.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var decodedBytes global.QManDBStruct
		err := rlp.Decode(bytes.NewReader(value), &decodedBytes)
		if err != nil {
			log.Info("Qmanager RLP", "Decoding Error", err.Error())
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		nowtimestamp, _ := time.Parse("2006-01-02 15:04:05", timestamp)
		dbtimestamp, _ := time.Parse("2006-01-02 15:04:05", decodedBytes.Timestamp)

		diff := nowtimestamp.Sub(dbtimestamp).Seconds()
		if diff > 30 {
			log.Info("Qman Expired Node", "Current Time", nowtimestamp)
			log.Info("Qman Expired Node", "DB Last Updated", dbtimestamp)

			log.Info("Qman Expired Node", "Expired Node", string(key))

			//TODO-REAP: Remove deleting node info due to stopping consensus issue.
			//err = global.QManagerStorage.Delete(key, nil)
			GetDBData()
			// if err != nil {
			// 	log.Info("Qmanager RLP", "Decoding Error", err.Error())
			// }
		}
	}
	CloseDB()
	//log.Info("Qmanager", "DB Status", "1. Disconnected")
}

func UpdateSenatorCandidateNodes() {
	//log.Info("Qmanager", "DB Status", "2. Connected")
	log.Info("Sentor & Candidate Update ", "List", global.GovernanceList)

	ConnectDB()
	for _, element := range global.GovernanceList {
		nodeAddress := element.Validator
		//log.Info("Sentor & Candidate Update ", "Node Address", nodeAddress)

		node_address_encoded, _ := rlp.EncodeToBytes(nodeAddress)
		foundNode, err := global.QManagerStorage.Get(node_address_encoded, nil)
		if err != nil {
			log.Info("Sentor & Candidate Update ", "DB Error", "Node Not Found")

		}
		if foundNode != nil {
			var decodedBytes global.QManDBStruct
			DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
			if DecodeErr != nil {
				log.Info("Sentor & Candidate Update", "Decoding Error", DecodeErr.Error())
			}
			var encodedStruct *global.QManDBStruct
			initBytes, err := rlp.EncodeToBytes(encodedStruct)

			if err != nil {
				log.Info("Sentor & Candidate Update", "RLP Error", err)
			}
			convertedTag := convertTagToString(element.Tag)
			//log.Info("Sentor & Candidate Update ", "Node Tag", convertedTag)
			encodedStruct = &global.QManDBStruct{ID: decodedBytes.ID, Address: decodedBytes.Address,
				Timestamp: decodedBytes.Timestamp, Tag: convertedTag}
			initBytes, err = rlp.EncodeToBytes(encodedStruct)
			if err != nil {
				log.Info("Sentor & Candidate Update", "RLP Error", err)
			}
			saveError := global.QManagerStorage.Put(node_address_encoded, initBytes, nil)

			if saveError != nil {
				log.Info("Sentor & Candidate Update", "err --", saveError)
			}
			GetDBData()

		}
	}
	CloseDB()
	//log.Info("Qmanager", "DB Status", "2. Disconnected")
}

func FindNode(dbStruct global.QManDBStruct) (found bool) {
	nodeAddress := dbStruct.Address
	node_address_encoded, err := rlp.EncodeToBytes(nodeAddress)
	if err != nil {
		log.Warn("QManager FindNode failed", "err", err)
		return false
	}
	//ecies.Encrypt()
	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	//log.Info("Qmanager", "DB Status", "3. Connected")
	ConnectDB()
	foundNode, err := global.QManagerStorage.Get(node_address_encoded, nil)
	if err != nil {
		//log.Info("Qmanager", "DB Status", "3. Disconnected")
		log.Info("QManager ", "DB --", "Node Not Found")
		CloseDB()
		return false
	}

	var decodedBytes global.QManDBStruct
	DecodeErr := rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
	if DecodeErr != nil {
		log.Info("Qmanager ", "Decoding Error", DecodeErr.Error())
		CloseDB()
		return false
	}

	var encodedStruct *global.QManDBStruct
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// encodedStruct = &global.QManDBStruct{ID: decodedBytes.ID, Address: decodedBytes.Address, Timestamp: timestamp, Tag: decodedBytes.Tag}
	encodedStruct = &global.QManDBStruct{ID: dbStruct.ID, Address: decodedBytes.Address, Timestamp: timestamp, Tag: decodedBytes.Tag}
	initBytes, err := rlp.EncodeToBytes(encodedStruct)
	if err != nil {
		log.Info("QManager DB Save", "err --", err)
		CloseDB()
		return false
	}
	updateError := global.QManagerStorage.Put(node_address_encoded, initBytes, nil)
	if updateError != nil {
		log.Info("QManager DB Save", "UPDATE ERROR --", updateError)
		CloseDB()
		return false
	}
	GetDBData()
	CloseDB()
	//log.Info("Qmanager", "DB Status", "3. Disconnected")

	return true
}

func Save(dbStruct global.QManDBStruct) (saved bool) {
	encodedAddress, _ := rlp.EncodeToBytes(dbStruct.Address)
	var encodedStruct *global.QManDBStruct
	initBytes, err := rlp.EncodeToBytes(encodedStruct)
	if err != nil {
		log.Info("QManager DB Save", "err --", err)
	}
	nodeTag := dbStruct.Tag
	for _, element := range global.GovernanceList {
		nodeAddress := element.Validator
		if dbStruct.Address == nodeAddress {
			nodeTag = convertTagToString(element.Tag)
			break
		}
	}
	encodedStruct = &global.QManDBStruct{ID: dbStruct.ID, Address: dbStruct.Address, Timestamp: dbStruct.Timestamp, Tag: nodeTag}
	initBytes, err = rlp.EncodeToBytes(encodedStruct)
	if err != nil {
		log.Info("QManager DB Save", "err --", err)
	}
	saved = SaveToDB(encodedAddress, initBytes)
	if !saved {
		return false
	}
	return true
}

func SaveToDB(Address []byte, NodeDetails []byte) (saved bool) {
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateRandomNumbers() (RandomNumbers []uint64, err error) {
	generatedBUFFER, err := quantum.GenerateQRNGData()
	if err == nil {
		one := make([]byte, 8)
		var RandomNums []uint64
		counter := 0
		for i := 0; i < len(generatedBUFFER)-4; i++ {
			one[counter] = generatedBUFFER[i]
			counter = counter + 1
			if counter == 4 {
				RandomNums = append(RandomNums, binary.LittleEndian.Uint64(one))
				counter = 0
			}
		}
		return RandomNums, nil
	} else {
		return nil, err
	}
}

func convertTagToString(typeTag common.Tag) (tag string) {
	if typeTag == common.Senator {
		return "0"
	} else if typeTag == common.Candidate {
		return "2"
	} else {
		return "3"
	}
}
