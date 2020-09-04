package qmanager

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var (
	Counter int
	Divisor int
	DBName string

	ConfigValidatorsParsed bool
)




func GovernanceSendList(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var govStruct []global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		m := global.Message{
			Message: "Error",
			Code: http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(m)
		return
	}
	global.GovernanceList = govStruct

	m := global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
	go  UpdateSenatorCandidateNodes()
}

func Ping(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	m := global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(m)
}


func Start(Addr *string, qmanKey *ecdsa.PrivateKey) {

	http.HandleFunc("/Ping", Ping)
	http.HandleFunc("/GovernanceSendList", GovernanceSendList)
	http.HandleFunc("/ExtraData", handleExtraData)
	http.HandleFunc("/BootNodeSendData", BootNodeSendData)
	http.HandleFunc("/CoordinatorConfirmation", CoordinatorConfirmation)

	s := strings.Split(*Addr, ":")
	DBName = s[1]
	http.ListenAndServe(*Addr, nil)

}

func  CoordinatorConfirmation(w http.ResponseWriter, req *http.Request)  {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Info("COORDINATOR CONFIRMATION")

	log.Info(string(body))

	var coordiStruct global.RequestCoordiStruct
	err = json.Unmarshal(body, &coordiStruct)
	if err != nil {
		panic(err)
	}
	log.Info("QMAN ", "DIVISOR: ", Divisor)

	if coordiStruct.QRND%uint64(Divisor) == 0 {
		log.Info("QMAN COORDI TRUE")

		decideStruct  := global.CoordiDecideStruct{
			Status: true,
		}
		json.NewEncoder(w).Encode(decideStruct)


	} else{
		log.Info("QMAN COORDI FALSE")


		decideStruct  := global.CoordiDecideStruct{
			Status: false,
		}
		json.NewEncoder(w).Encode(decideStruct)
	}


}

func  BootNodeSendData (w http.ResponseWriter, req *http.Request){

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Info("Qmanager is not alive")
		panic(err)
	}
	var nodeStruct global.QManDBStruct
	err = json.Unmarshal(body, &nodeStruct)
	if err != nil {
		panic(err)
	}

	log.Info("Bootnode Data ", "Addr: ", nodeStruct.Address )

	if nodeStruct.Address != ""{
		if !FindNode(nodeStruct.Address){
			Save(nodeStruct)
		}
	}

	m := global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)

}

func  handleExtraData (w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		log.Info(string(body))

		var reqStruct global.RequestStruct
		err = json.Unmarshal(body, &reqStruct)
		if err != nil {
		panic(err)
	}

		proposerAddress := common.HexToAddress(reqStruct.Proposer)

		log.Info("Received EXTRA DATA REQUEST from geth")
		if global.QRNGDeviceStat == true{
			log.Info("Random Number Generator " ,  "Using - " , "Quantum Device")
		} else {
			log.Info("Random Number Generator " ,  "Using - " , "Pusedo Random")
		}

		Counter = Counter + 1
		log.Info("Round ", "Count: ", Counter)

		var extra []common.ValidatorInfo

		outerLoop := 0
		for {
			//log.Print("Qmanager ", "Generating Random Numbers ", "Outerloop")
			extra = generateExtraData()
			completed := false
			divisor := rand.Intn(50) + 1

			index := 0
			for index < len(extra) {

				if  proposerAddress != extra[index].Address {
					if extra[index].Tag == common.Senator{
						randomNumber := extra[index].Qrnd
						if randomNumber%uint64(divisor) == 0 {
							extra[index].Tag = common.Coordinator
							log.Info("Qmanager " , "Random Coordinator Selected ", extra[index].Address.String())
							index = len(extra)
							completed = true
							Divisor = divisor
						}
					}
				}
				index++
			}
			outerLoop++
			if completed{
				log.Info("QManager ExtraData ", "For Loop Index: ", outerLoop)
				break
			}
			if outerLoop == 30{
				log.Error("QManager ExtraData ", "Error", "Cannot Select Coordinator")
				break
			}

		}
		log.Info("QManager ", "ExtraData Length: ", len(extra))
		log.Info("QManager ", "ExtraData: ", extra)
		json.NewEncoder(w).Encode(extra)
}


func generateExtraData() []common.ValidatorInfo{

	var extra []common.ValidatorInfo
	for _, validator := range global.DBDataList {

		var num uint64

		if global.QRNGDeviceStat == true{
 			randomIndex := rand.Intn(12280) + 1
			num = global.RandomNumbers[randomIndex]

		} else {
			num = rand.Uint64()
		}
		validatorInfo := common.ValidatorInfo{}
		validatorInfo.Address = common.HexToAddress(validator.Address)
		validatorInfo.Qrnd = num
		convertedTag, _ := strconv.Atoi(validator.Tag)
		validatorInfo.Tag = common.Tag(convertedTag)

		extra = append(extra, validatorInfo)
	}
	return extra
}
