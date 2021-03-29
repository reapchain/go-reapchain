package qmanager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	Counter int
	//Divisor int
	CoordiQrnd uint64
	DBName     string

	ConfigValidatorsParsed bool
)

func Start(Addr *string, qmanKey *ecdsa.PrivateKey) {

	http.HandleFunc("/Ping", Ping)
	http.HandleFunc("/GovernanceSendList", GovernanceSendList)
	http.HandleFunc("/GovernanceAddValidators", GovernanceAddValidators)
	http.HandleFunc("/GovernanceRemoveValidators", GovernanceRemoveValidators)
	http.HandleFunc("/GovernanceGetValidatorList", GovernanceGetValidatorList)
	http.HandleFunc("/ExtraData", handleExtraData)
	http.HandleFunc("/BootNodeSendData", BootNodeSendData)
	http.HandleFunc("/CoordinatorConfirmation", CoordinatorConfirmation)

	s := strings.Split(*Addr, ":")
	DBName = s[1]
	http.ListenAndServe(*Addr, nil)

}

func Ping(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}

func GovernanceSendList(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		//panic(err)
		return
	}

	var govStruct []global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		m := global.Message{
			Message: "Error",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(m)
		return
	}
	global.GovernanceList = govStruct

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
	go UpdateSenatorCandidateNodes()
}

func GovernanceAddValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		return
	}

	var govStruct []global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		m := global.Message{
			Message: "Error",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	ConnectDB()
	list := make(map[string]common.Tag, len(govStruct))
	for _, node := range govStruct {
		if node.Tag != common.Senator && node.Tag != common.Candidate {
			continue
		}

		encodekey, _ := rlp.EncodeToBytes(node.Validator)
		foundNode, err := global.QManagerStorage.Get(encodekey, nil)
		if err != nil {
			log.Debug("Node not found in DB", "error", err)
			if err != leveldb.ErrNotFound {
				continue
			}
		}

		id := ""
		if foundNode != nil {
			var decodedBytes global.QManDBStruct
			err = rlp.Decode(bytes.NewReader(foundNode), &decodedBytes)
			if err != nil {
				log.Warn("Rlp decode error", "error", err)
				continue
			}
			id = decodedBytes.ID
		}

		var encodedStruct *global.QManDBStruct
		//convertedTag := convertTagToString(node.Tag)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		//encodedStruct = &global.QManDBStruct{ID: decodedBytes.ID, Address: decodedBytes.Address, Timestamp: timestamp, Tag: convertTagToString(node.Tag)}
		encodedStruct = &global.QManDBStruct{ID: id, Address: node.Validator, Timestamp: timestamp, Tag: convertTagToString(node.Tag)}
		initBytes, err := rlp.EncodeToBytes(encodedStruct)
		if err != nil {
			log.Warn("Rlp encode error", "error", err)
			continue
		}

		err = global.QManagerStorage.Put(encodekey, initBytes, nil)
		if err != nil {
			log.Warn("DB Save error", "error", err)
			continue
		}

		list[node.Validator] = node.Tag
	}
	GetDBData()
	CloseDB()

	for i, vlist := range global.GovernanceList {
		if tag, ok := list[vlist.Validator]; ok {
			global.GovernanceList[i].Tag = tag
		} else {
			global.GovernanceList = append(global.GovernanceList, vlist)
		}
	}

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}

func GovernanceRemoveValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		return
	}

	var govStruct []global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		m := global.Message{
			Message: "Error",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	ConnectDB()
	list := make(map[string]common.Tag, len(govStruct))
	for _, node := range govStruct {
		encodekey, _ := rlp.EncodeToBytes(node.Validator)
		err = global.QManagerStorage.Delete(encodekey, nil)
		if err != nil {
			log.Warn("DB delete error", "error", err)
			continue
		}

		list[node.Validator] = node.Tag
	}
	GetDBData()
	CloseDB()

	for i, vlist := range global.GovernanceList {
		if _, ok := list[vlist.Validator]; ok {
			global.GovernanceList = append(global.GovernanceList[:i], global.GovernanceList[i+1:]...)
		}
	}

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}

func GovernanceGetValidatorList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		return
	}

	var addrs []string
	err = json.Unmarshal(body, &addrs)
	if err != nil {
		log.Error("request body error", "err", err)
		m := global.Message{
			Message: "Error",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	ConnectDB()
	GetDBData()
	CloseDB()

	var retList []global.QManDBStruct
	if len(addrs) != 0 {
		list := make(map[string]int, len(addrs))
		for i, addr := range addrs {
			list[addr] = i
		}
		for _, vlist := range global.DBDataList {
			if _, ok := list[vlist.Address]; ok {
				retList = append(retList, vlist)
			}
		}
	} else {
		retList = global.DBDataList
	}

	json.NewEncoder(w).Encode(retList)
}

func CoordinatorConfirmation(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		//panic(err)
		return
	}
	log.Info("COORDINATOR CONFIRMATION")

	log.Info(string(body))

	var coordiStruct global.RequestCoordiStruct
	err = json.Unmarshal(body, &coordiStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		//panic(err)
		return
	}
	// log.Info("QMAN ", "DIVISOR: ", Divisor)
	log.Info("QMAN ", "CoordiQrnd: ", CoordiQrnd)

	// if coordiStruct.QRND%uint64(Divisor) == 0 {
	if coordiStruct.QRND == CoordiQrnd {
		log.Info("QMAN COORDI TRUE")

		decideStruct := global.CoordiDecideStruct{
			Status: true,
		}
		json.NewEncoder(w).Encode(decideStruct)

	} else {
		log.Info("QMAN COORDI FALSE")

		decideStruct := global.CoordiDecideStruct{
			Status: false,
		}
		json.NewEncoder(w).Encode(decideStruct)
	}
}

func BootNodeSendData(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		//panic(err)
		return
	}
	var nodeStruct global.QManDBStruct
	err = json.Unmarshal(body, &nodeStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		//panic(err)
		return
	}

	log.Info("Bootnode Data ", "Addr: ", nodeStruct.Address)

	if nodeStruct.Address != "" {
		if !FindNode(nodeStruct) {
			Save(nodeStruct)
		}
	}

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)
}

func handleExtraData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("request failed", "err", err)
		//panic(err)
		return
	}
	log.Info(string(body))

	var reqStruct global.RequestStruct
	err = json.Unmarshal(body, &reqStruct)
	if err != nil {
		log.Error("request body error", "err", err)
		//panic(err)
		return
	}

	proposerAddress := common.HexToAddress(reqStruct.Proposer)

	log.Info("Received EXTRA DATA REQUEST from geth", "proposer", proposerAddress.Hex())
	if global.QRNGDeviceStat == true {
		log.Info("Random Number Generator ", "Using - ", "Quantum Device")
	} else {
		log.Info("Random Number Generator ", "Using - ", "Pusedo Random")
	}

	Counter = Counter + 1
	log.Info("Round ", "Count: ", Counter)

	var extra []common.ValidatorInfo

	outerLoop := 0
	for {
		//log.Print("Qmanager ", "Generating Random Numbers ", "Outerloop")
		extra = generateExtraData()
		completed := false
		if len(extra) == 0 {
			log.Error("There is no validator")
			break
		}
		//divisor := rand.Intn(50) + 1
		divisor := rand.Intn(50)%len(extra) + 1 //TODO-REAP: Need to check algorithm

		index := 0
		for index < len(extra) {
			if proposerAddress != extra[index].Address {
				if extra[index].Tag == common.Senator {
					randomNumber := extra[index].Qrnd
					if randomNumber%uint64(divisor) == 0 {
						extra[index].Tag = common.Coordinator
						log.Info("Qmanager ", "Random Coordinator Selected ", extra[index].Address.String())
						index = len(extra)
						completed = true
						//Divisor = divisor
						CoordiQrnd = randomNumber
					}
				}
			}
			index++
		}
		outerLoop++
		if completed {
			log.Info("QManager ExtraData ", "For Loop Index: ", outerLoop)
			break
		}
		if outerLoop == 30 {
			log.Error("QManager ExtraData ", "Error", "Cannot Select Coordinator")
			break
		}

	}
	log.Info("QManager ", "ExtraData Length: ", len(extra))
	log.Info("QManager ", "ExtraData: ", extra)
	json.NewEncoder(w).Encode(extra)
}

func generateExtraData() []common.ValidatorInfo {
	var extra []common.ValidatorInfo
	for _, validator := range global.DBDataList {
		nodeTag, _ := strconv.Atoi(validator.Tag)
		if common.Tag(nodeTag) != common.Senator && common.Tag(nodeTag) != common.Candidate {
			log.Debug("Normal node", "addr", validator.Address)
			continue
		}
		if validator.ID == "" {
			log.Debug("Node is not alive", "addr", validator.Address)
			continue
		}

		var num uint64

		if global.QRNGDeviceStat == true {
			randomIndex := rand.Intn(12280) + 1
			num = global.RandomNumbers[randomIndex]
		} else {
			num = rand.Uint64()
		}
		validatorInfo := common.ValidatorInfo{}
		validatorInfo.Address = common.HexToAddress(validator.Address)
		validatorInfo.Qrnd = num
		//convertedTag, _ := strconv.Atoi(validator.Tag)
		validatorInfo.Tag = common.Tag(nodeTag)

		extra = append(extra, validatorInfo)
	}
	return extra
}
