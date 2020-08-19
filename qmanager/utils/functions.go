package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"net"
	"net/http"
	"strings"
)

var (
	ActiveQmanager string
)

func CheckQmanagerStatus() bool {
	var QManagerAddresses= config.Config.QManagers

	for _, qman := range QManagerAddresses{
		split := strings.Split(qman, "@")
		QManager := split[1]
		ActiveQmanager = QManager
		//timeout := 1 * time.Second
		//conn, err := net.DialTimeout("http", qman + "/Ping" , timeout)
		//if err != nil {
		//	log.Info("QManager Not Available", "ADDR", QManager )
		//} else {
		//	conn.Close()
		//	ActiveQmanager = QManager
		//}
	}
	split := strings.Split(ActiveQmanager, ":")
	Address := split[0]
	trial := net.ParseIP(Address)
	if trial == nil {
		log.Error("QManager Adrress Error", "Invalid Address ", trial )
		return false
	} else {
		return true
	}

}

func CheckAddressValidity() bool {
	var QManagerAddresses= config.Config.QManagers
	if  len(QManagerAddresses) == 0 {
			log.Error("QManager Connection Error", "Address Not Found ", "Please insert qmanager address into Config.json" )
		return false
	} else{
		if CheckQmanagerStatus(){
			return true
		}else {
			return false
		}
	}
}

func RequestExtraData(Proposer string) (common.ValidatorInfos, error) {
	if CheckAddressValidity(){

		requestStruct := global.RequestStruct{
			Proposer: Proposer,
		}

		bytesRepresentation, err := json.Marshal(requestStruct)
		if err != nil {
			log.Error(err.Error())
		}

		log.Info("Get ExtraData", "QMANAGER Address : ", "http://"+ ActiveQmanager + "/ExtraData")

		resp, err := http.Post("http://"+ ActiveQmanager + "/ExtraData", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		var result []common.ValidatorInfo
		json.NewDecoder(resp.Body).Decode(&result)
		log.Info("VALIDATOR LIST", "Full List : ", result)
		//log.Info("VALIDATOR LIST", "Full BODY : ", resp)

		return result, nil

	}else {
		return nil, errors.New("Unavailable QManager Address")
	}
}

func BootNodeToQmanager(NodeData global.QManDBStruct) error {
	if CheckAddressValidity(){
		bytesRepresentation, err := json.Marshal(NodeData)
		if err != nil {
			log.Error(err.Error())
		}
		resp, err := http.Post("http://"+ ActiveQmanager + "/BootNodeSendData", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		var result global.Message

		json.NewDecoder(resp.Body).Decode(&result)

		log.Info("Bootnode To Qmanager", "Send Status : ", result.Message)
		return nil
	}else {
		return errors.New("Unavailable QManager Address")
	}
}

func CoordinatorConfirmation(coordiReq global.RequestCoordiStruct) (bool, error) {
	if CheckAddressValidity(){
	bytesRepresentation, err := json.Marshal(coordiReq)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	resp, err := http.Post("http://"+ ActiveQmanager + "/CoordinatorConfirmation", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	var result global.CoordiDecideStruct

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Status, nil

	}else {
		return false, errors.New("Unavailable QManager Address")
	}

}


