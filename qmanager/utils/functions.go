package utils

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qmanager/global"
	"net/http"
	"strings"
)

func RequestExtraData(Properser string) common.ValidatorInfos {
	if  len(config.Config.QManagers) == 0 {
		log.Info("File Not Found:", "QManagers", config.Config.QManagers[0] )
	}
	var QManagerURLs= config.Config.QManagers[0]
	s := strings.Split(QManagerURLs, "@")
	QManagerURL := s[1]

	requestStruct := global.RequestStruct{
		Proposer: Properser,
	}
	//message := map[string]interface{}{
	//	"hello": "world",
	//	"life":  42,
	//	"embedded": map[string]string{
	//		"yes": "of course!",
	//	},
	//}

	bytesRepresentation, err := json.Marshal(requestStruct)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("GET EXTRADAT", "QMANAGER URL : ", "http://"+ QManagerURL + "/ExtraData")

	resp, err := http.Post("http://"+ QManagerURL + "/ExtraData", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Error(err.Error())
	}


	//var reqStruct []common.ValidatorInfo
	//if err != nil {
	//	panic(err)
	//}

	var result []common.ValidatorInfo
	json.NewDecoder(resp.Body).Decode(&result)
	log.Info("VALIDATOR LIST", "Full List : ", result)

	log.Info("VALIDATOR LIST", "Full BODY : ", resp)

	//data := []common.ValidatorInfo{}
	//json.Unmarshal([]byte(s), &data)

	return result

	//log.Info(result["data"])
}

func BootNodeToQmanager(NodeData global.QManDBStruct) {
    if  len(config.Config.QManagers) == 0 {
    	log.Info("File Not Found:", "QManagers", config.Config.QManagers[0] )
	}

	var QManagerURLs = config.Config.QManagers[0]
	s := strings.Split(QManagerURLs, "@")
	QManagerURL := s[1]

	//message := map[string]interface{}{
	//	"hello": "world",
	//	"life":  42,
	//	"embedded": map[string]string{
	//		"yes": "of course!",
	//	},
	//}

	bytesRepresentation, err := json.Marshal(NodeData)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post("http://"+ QManagerURL + "/BootNodeSendData", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Error(err.Error())
	}

	//json.NewDecoder(resp.Body).Decode(&result)

	var result global.Message

	json.NewDecoder(resp.Body).Decode(&result)

	log.Info("Bootnode To Qmanager", "Send Status : ", result.Message)


}

func CooridnatorConfirmation(coordiReq global.RequestCoordiStruct) bool {
	if  len(config.Config.QManagers) == 0 {
		log.Info("File Not Found:", "QManagers", config.Config.QManagers[0] )
	}
	var QManagerURLs= config.Config.QManagers[0]
	s := strings.Split(QManagerURLs, "@")
	QManagerURL := s[1]

	//message := map[string]interface{}{
	//	"hello": "world",
	//	"life":  42,
	//	"embedded": map[string]string{
	//		"yes": "of course!",
	//	},
	//}

	bytesRepresentation, err := json.Marshal(coordiReq)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	resp, err := http.Post("http://"+ QManagerURL + "/CoordinatorConfirmation", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Error(err.Error())
		return false
	}




	//json.NewDecoder(resp.Body).Decode(&result)


	var result global.CoordiDecideStruct

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Status



}


func CheckBootNodePortAndID(NodeID string, Port int) bool{
	//log.Info("BootNode", "ID = ", NodeID)
	//log.Info("BootNode", "IP Addr = ", Port)
	//
	//log.Info("Global", "ID = ", BootNodeID)
	//log.Info("Global", "IP Addr = ", BootNodePort)

	if NodeID == BootNodeID && Port == BootNodePort {
		return true
	}

	return false
}