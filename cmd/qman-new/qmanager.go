package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager"
	"github.com/ethereum/go-ethereum/qmanager/global"
)

type roundInfo struct {
	// coordinator common.Address
	// cQrnd       int
	divisor    int
	validators []common.ValidatorInfo
}

type Qmanager struct {
	// seq   uint64
	// round uint64
	addr string
	// port  uint64
	// dir   string

	all map[common.Address]common.Tag
	// senators   map[common.Address]int
	// candidates map[common.Address]int
	//alive      map[common.Address]int
	current roundInfo

	config *config.EnvConfig
	db     ethdb.Database // Block chain database
	//p2pServer *p2p.Server
	//server *http.Server

	lock sync.RWMutex
}

func NewQmanager(db ethdb.Database, config *config.EnvConfig, addr string) *Qmanager {
	log.Info("NewQmanager")
	qman := &Qmanager{
		config: config,
		db:     db,
		addr:   addr,
		all:    make(map[common.Address]common.Tag),
	}

	//load latest validators from db
	senators := ReadSenators(qman.db)
	if len(senators) != 0 {
		for i, senator := range senators {
			qman.all[senator] = common.Senator
			log.Debug("db senator", "i", i, "addr", senator.Hex())
		}
	} else {
		for _, senator := range config.Senatornodes {
			qman.all[common.HexToAddress(senator)] = common.Senator
			if err := WriteSenator(qman.db, common.HexToAddress(senator), int(common.Senator)); err != nil {
				continue
			}
			log.Debug("config.json senator", "addr", senator, "value", common.Senator)
		}
	}
	candidates := ReadCandidates(qman.db)
	if len(candidates) != 0 {
		for i, candidate := range candidates {
			qman.all[candidate] = common.Candidate
			log.Debug("db candidate", "i", i, "addr", candidate.Hex())
		}
	} else {
		for _, candidate := range config.Candidatenodes {
			qman.all[common.HexToAddress(candidate)] = common.Candidate
			if err := WriteCandidate(qman.db, common.HexToAddress(candidate), int(common.Candidate)); err != nil {
				continue
			}
			log.Debug("config.json candidate", "addr", candidate, "value", common.Candidate)
		}
	}

	qman.printInfo()

	qmanager.CheckQRNGStatus()

	http.HandleFunc("/Ping", qmanager.Ping)
	//http.HandleFunc("/GovernanceSendList", qman.GovernanceSendList)
	http.HandleFunc("/GovernanceAddValidators", qman.AddValidators)
	http.HandleFunc("/GovernanceRemoveValidators", qman.RemoveValidators)
	http.HandleFunc("/GovernanceGetValidatorList", qman.GetValidatorList)
	http.HandleFunc("/ExtraData", qman.HandleExtraData)
	http.HandleFunc("/BootNodeSendData", qman.BootNodeSendData)
	http.HandleFunc("/CoordinatorConfirmation", qman.CoordinatorConfirmation)

	return qman
}

func (qm *Qmanager) Start() {
	http.ListenAndServe(qm.addr, nil)
	log.Info("Qmanager started")
}

// func (qm *Qmanager) GovernanceSendList(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	body, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var govStruct []global.GovStruct
// 	err = json.Unmarshal(body, &govStruct)
// 	if err != nil {
// 		m := global.Message{
// 			Message: "Error",
// 			Code:    http.StatusBadRequest,
// 		}
// 		json.NewEncoder(w).Encode(m)
// 		return
// 	}
// 	global.GovernanceList = govStruct
// 	m := global.Message{
// 		Message: "Success",
// 		Code:    http.StatusOK,
// 	}
// 	json.NewEncoder(w).Encode(m)
// 	go qmanager.UpdateSenatorCandidateNodes()
// }

func (qm *Qmanager) AddValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	var validators []global.GovStruct
	err = json.Unmarshal(body, &validators)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	for _, validator := range validators {
		if validator.Tag != common.Senator && validator.Tag != common.Candidate {
			log.Warn("invalid tag", "addr", validator.Validator, "tag", validator.Tag)
			continue
		}
		if validator.Tag == common.Senator {
			if err := WriteSenator(qm.db, common.HexToAddress(validator.Validator), int(validator.Tag)); err != nil {
				continue
			}
		} else if validator.Tag == common.Candidate {
			if err := WriteCandidate(qm.db, common.HexToAddress(validator.Validator), int(validator.Tag)); err != nil {
				continue
			}
		}
		qm.all[common.HexToAddress(validator.Validator)] = validator.Tag
		log.Info("add validator", "addr", validator.Validator, "tag", validator.Tag)
	}

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)

	qm.printInfo()
}

func (qm *Qmanager) RemoveValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	var validators []global.GovStruct
	err = json.Unmarshal(body, &validators)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	for _, validator := range validators {
		if validator.Tag != common.Senator && validator.Tag != common.Candidate {
			log.Warn("invalid tag", "addr", validator.Validator, "tag", validator.Tag)
			continue
		}
		if validator.Tag == common.Senator {
			if err := RemoveSenator(qm.db, common.HexToAddress(validator.Validator)); err != nil {
				continue
			}
		} else if validator.Tag == common.Candidate {
			if err := RemoveCandidate(qm.db, common.HexToAddress(validator.Validator)); err != nil {
				continue
			}
		}
		delete(qm.all, common.HexToAddress(validator.Validator))
		log.Info("remove validator", "addr", validator.Validator, "tag", validator.Tag)
	}

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)

	qm.printInfo()
}

func (qm *Qmanager) GetValidatorList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}

	var addrs []string
	err = json.Unmarshal(body, &addrs)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}

	// var retList []global.QManDBStruct
	var retList []global.GovStruct
	if len(addrs) != 0 {
		for _, addr := range addrs {
			if tag, ok := qm.all[common.HexToAddress(addr)]; ok {
				retList = append(retList, global.GovStruct{Validator: addr, Tag: tag})
			} else {
				//read db, if success, return data and update qm.all
			}
		}
	} else {
		for addr, tag := range qm.all {
			retList = append(retList, global.GovStruct{Validator: addr.Hex(), Tag: tag})
		}
	}
	log.Info("Get validator list", "count", len(retList))

	//qm.responseSuccess(w)
	//json.NewEncoder(w).Encode(retList)
	qm.responseData(w, retList)

	qm.printInfo()
}

func (qm *Qmanager) Ping(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}

func (qm *Qmanager) HandleExtraData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	log.Info(string(body))

	var reqStruct global.RequestStruct
	err = json.Unmarshal(body, &reqStruct)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}

	proposerAddress := common.HexToAddress(reqStruct.Proposer)

	log.Info("Received EXTRA DATA REQUEST from geth")
	if global.QRNGDeviceStat == true {
		log.Info("Random Number Generator ", "Using - ", "Quantum Device")
	} else {
		log.Info("Random Number Generator ", "Using - ", "Pusedo Random")
	}

	var extra []common.ValidatorInfo

	outerLoop := 0
	for {
		//log.Print("Qmanager ", "Generating Random Numbers ", "Outerloop")
		extra = qm.generateExtraData()
		completed := false
		divisor := rand.Intn(50) + 1

		index := 0
		log.Debug("handleExtraData", "len(extra)", len(extra))
		for index < len(extra) {
			if proposerAddress != extra[index].Address {
				if extra[index].Tag == common.Senator {
					randomNumber := extra[index].Qrnd
					if randomNumber%uint64(divisor) == 0 {
						extra[index].Tag = common.Coordinator
						log.Debug("Qmanager ", "Random Coordinator Selected ", extra[index].Address.String())
						index = len(extra)
						completed = true
						qm.current = roundInfo{
							divisor:    divisor,
							validators: extra,
						}
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

func (qm *Qmanager) generateExtraData() []common.ValidatorInfo {
	var extra []common.ValidatorInfo
	for addr, tag := range qm.all {
		log.Debug("generateExtraData", "validator.Address", addr, "tag", tag)

		var num uint64
		if global.QRNGDeviceStat == true {
			randomIndex := rand.Intn(12280) + 1
			num = global.RandomNumbers[randomIndex]
		} else {
			num = rand.Uint64()
		}
		validatorInfo := common.ValidatorInfo{
			Address: addr,
			Qrnd:    num,
			Tag:     tag,
		}
		extra = append(extra, validatorInfo)
	}
	return extra
}

func (qm *Qmanager) CoordinatorConfirmation(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	log.Info("COORDINATOR CONFIRMATION")
	log.Debug(string(body))

	var coordiStruct global.RequestCoordiStruct
	err = json.Unmarshal(body, &coordiStruct)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	log.Info("QMAN ", "DIVISOR: ", qm.current.divisor)

	if coordiStruct.QRND%uint64(qm.current.divisor) == 0 {
		log.Info("QMAN COORDI TRUE")
		decideStruct := global.CoordiDecideStruct{Status: true}
		json.NewEncoder(w).Encode(decideStruct)
	} else {
		log.Info("QMAN COORDI FALSE")
		decideStruct := global.CoordiDecideStruct{Status: false}
		json.NewEncoder(w).Encode(decideStruct)
	}
}

func (qm *Qmanager) BootNodeSendData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warn("http request fail", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}
	var nodeStruct global.QManDBStruct
	err = json.Unmarshal(body, &nodeStruct)
	if err != nil {
		log.Warn("json unmarshal error", "err", err)
		qm.responseFail(w, http.StatusBadRequest)
		return
	}

	log.Info("Bootnode Data ", "Addr: ", nodeStruct.Address, "Tag", nodeStruct.Tag)

	// if nodeStruct.Address != "" {
	// 	if !FindNode(nodeStruct.Address) {
	// 		Save(nodeStruct)
	// 	}
	// }

	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)
}

func (qm *Qmanager) printInfo() {
	for addr, tag := range qm.all {
		log.Debug("print all node infos", "addr", addr.Hex(), "tag", tag)
	}
}

// func (qm *Qmanager) response(w http.ResponseWriter, msg string, code int) {
// 	m := global.Message{
// 		Message: msg,
// 		Code:    code,
// 	}
// 	json.NewEncoder(w).Encode(m)
// }

func (qm *Qmanager) responseFail(w http.ResponseWriter, code int) {
	m := global.Message{
		Message: "Fail",
		Code:    code,
	}
	json.NewEncoder(w).Encode(m)
}

func (qm *Qmanager) responseSuccess(w http.ResponseWriter) {
	m := global.Message{
		Message: "Success",
		Code:    http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}

func (qm *Qmanager) responseData(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
