package podc_global

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
)

func RequestExtraData(Properser string) common.ValidatorInfos {

	var QManagerURL= config.Config.QManagers[0]

	requestStruct := RequestStruct{
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

	return result

	//log.Info(result["data"])
}

func BootNodeSendData(NodeData QManDBStruct) {

	var QManagerURL= config.Config.QManagers[0]


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


	var result Message

	json.NewDecoder(resp.Body).Decode(&result)

	log.Info("BootnodeToQmanager", "Send Status : ", result.Message)


}

func CooridnatorConfirmation(coordiReq RequestCoordiStruct) bool {

	var QManagerURL= config.Config.QManagers[0]


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

	resp, err := http.Post("http://"+ QManagerURL + "/BootNodeSendData", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Error(err.Error())
		return false
	}




	//json.NewDecoder(resp.Body).Decode(&result)


	var result CoordiDecideStruct

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Status



}
