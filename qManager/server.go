package qManager

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"strings"

	"io/ioutil"
	"math/rand"
	"net/http"

	"fmt"
)




var (
	Counter int
	Divisor int
	DBName string
)

//func ECCDecrypt(ct []byte, prk ecies.PrivateKey) ([]byte, error) {
//	pt, err := prk.Decrypt(rand.Reader, ct, nil, nil)
//	return pt, err
//}
//func hello(w http.ResponseWriter, req *http.Request) {
//
//	// fmt.Fprintf(w, "hello\n")
//	w.Header().Set("Content-Type", "application/json")
//	body, err := ioutil.ReadAll(req.Body)
//    if err != nil {
//        panic(err)
//    }
//	log.Println(string(body))
//    var t []test_struct
//    err = json.Unmarshal(body, &t)
//    if err != nil {
//        panic(err)
//    }
//	log.Println(t[0].Validator)
//	m := Message{
//		Message: "Success",
//		Code: http.StatusOK,
//	}
//	// b, err := json.Marshal(m)
//
//	// p := "[{'Code': 'SUCCESS'},]"
//	json.NewEncoder(w).Encode(m)
//	b, err := json.Marshal(m)
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//	// w.WriteHeader(http.StatusCreated)
//    // json.NewEncoder(w).Encode(data)
//}

func RequestQmanager(w http.ResponseWriter, req *http.Request) {

	// fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	//log.Print(string(body))

	//log.Info("QManager Server Started")


	var govStruct []podc_global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		m := podc_global.Message{
			Message: "Error",
			Code: http.StatusBadRequest,
		}

		json.NewEncoder(w).Encode(m)
		return
	}





	// // b, err := json.Marshal(m)

	// // p := "[{'Code': 'SUCCESS'},]"
	//for index, element := range t {
	//	// log.Println(element)
	//	// log.Println("%d",index);
	//	log.Printf("Index: %d, Address:  %s, Tag: %s", index, element.Validator, element.Tag)
	//	// log.Println("Index: " + index + ", Address: " + element.Validator + ", Tag: " + element.Tag)
	//	// index is the index where we are
	//	// element is the element from someSlice for where we are
	//}

	// privateKeyFile, err := os.Open("private_key.pem")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// pemfileinfo, _ := privateKeyFile.Stat()
	// var size int64 = pemfileinfo.Size()
	// pembytes := make([]byte, size)
	// buffer := bufio.NewReader(privateKeyFile)
	// _, err = buffer.Read(pembytes)
	// data, _ := pem.Decode([]byte(pembytes))
	// privateKeyFile.Close()

	// privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// var privKey *ecdsa.PrivateKey
	// var ok bool
	// privKey, ok = privateKeyImported.(*ecdsa.PrivateKey)
	// if !ok {
	// 	fmt.Println("Error")

	// }

	// fmt.Println(privKey)

	// eciesKey := ecies.ImportECDSA(privKey)
	// dedata, err := ECCDecrypt(body, *eciesKey)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Decrypt", string(dedata))
	//timestamp := time.Now().String()
	//ts, _ := time.Parse(timestamp, "2006-01-02 15:04:05")
	//fmt.Println(ts)


	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(data)

	podc_global.GovernanceList = govStruct

	m := podc_global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)

	go  UpdateSenatorCandidateNodes()

	//log.Printf("HTTP Server Response Sent")


}
func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}




func Start(Addr *string, qmanKey *ecdsa.PrivateKey) {



	// bodyBytes := []byte{91, 66, 64, 51, 102, 100, 97, 100, 49, 56, 56}

	//timestamp := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(timestamp)
	//
	//secondTime := time.Now().Add(time.Second * time.Duration(12)).Format("2006-01-02 15:04:05")
	//
	//fmt.Println(secondTime)
	//
	//
	//t, _ := time.Parse("2006-01-02 15:04:05", timestamp )
	//t2, _ := time.Parse("2006-01-02 15:04:05", secondTime )
	//fmt.Println(t)
	//
	//
	//
	//diff := t2.Sub(t)
	//fmt.Println(diff)
	// s := string(bodyBytes[:])
	// log.Println(s)
	// myString := hex.EncodeToString(bodyBytes)
	// log.Println(myString)
	//http.HandleFunc("/hello", hello)
	http.HandleFunc("/RequestQmanager", RequestQmanager)
	http.HandleFunc("/ExtraData", handleExtraData)
	http.HandleFunc("/BootNodeSendData", BootNodeSendData)
	http.HandleFunc("/CoordinatorConfirmation", CoordinatorConfirmation)



	//addr, err := net.ResolveUDPAddr("udp", *Addr)
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

	var coordiStruct podc_global.RequestCoordiStruct
	err = json.Unmarshal(body, &coordiStruct)
	if err != nil {
		panic(err)
	}


	log.Info("QMAN ", "DIVISOR: ", Divisor)


	if coordiStruct.QRND%uint64(Divisor) == 0 {
		log.Info("QMAN COORDI TRUE")

		decideStruct  := podc_global.CoordiDecideStruct{
			Status: true,
		}
		json.NewEncoder(w).Encode(decideStruct)


	} else{
		log.Info("QMAN COORDI FALSE")


		decideStruct  := podc_global.CoordiDecideStruct{
			Status: false,
		}
		json.NewEncoder(w).Encode(decideStruct)
	}


}

func  BootNodeSendData (w http.ResponseWriter, req *http.Request){

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	//log.Info(string(body))

	var nodeStruct podc_global.QManDBStruct
	err = json.Unmarshal(body, &nodeStruct)
	if err != nil {
		panic(err)
	}

	log.Info("Bootnode Data ", "Addr: ", nodeStruct.Address)


	//log.Info("Bootnode Sent Data Address" )
	//log.Info(nodeStruct.Address)


	if nodeStruct.Address != ""{
		if !FindNode(nodeStruct.Address){
			Save(nodeStruct)
		}
	}

	m := podc_global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)



}

//
////For Qmanager, event handler to receive msg from geth
func  handleExtraData (w http.ResponseWriter, req *http.Request){

	//ConnectDB()
	//GetDBData()
	//CloseDB()

		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		log.Info(string(body))

		var reqStruct podc_global.RequestStruct
		err = json.Unmarshal(body, &reqStruct)
		if err != nil {
		panic(err)
	}

		proposerAddress := common.HexToAddress(reqStruct.Proposer)

		log.Info("Received EXTRA DATA REQUEST from geth")

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
				//log.Print("Qmanager ", "Generating Random Numbers ", "InnerLoop")

				if  proposerAddress != extra[index].Address {
					//if extra[index].Tag == common.Senator{
						randomNumber := extra[index].Qrnd
						if randomNumber%uint64(divisor) == 0 {
							extra[index].Tag = common.Coordinator
							log.Info("Qmanager " , "Random Coordinator Selected ", extra[index].Address.String())
							index = len(extra)
							completed = true
							Divisor = divisor
						}
					//}
				}
				//log.Print("ExtraData list", "Address", extra[index].Address , "Qrnd", extra[index].Qrnd, "Tag",  extra[index].Tag)
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


		//log.Print("ExtraData list", "extradata", extra)

		//defer db.Close()
		log.Info("QManager ", "ExtraData Length: ", len(extra))
		log.Info("QManager ", "ExtraData: ", extra)
	//	extraDataJson, err := json.Marshal(extra)
	//
	//	if err != nil {
	//		log.Print("Failed to encode JSON", err)
	//	}
	//
	//log.Print("QManager", "ExtraDataJSON: ", extraDataJson)


		//m := podc_global.Message{
		//		Message: "Success",
		//		Code: http.StatusOK,
		//	}

		json.NewEncoder(w).Encode(extra)

}


func generateExtraData() []common.ValidatorInfo{

	//qManager.ConnectDB()
	//log.Info("Qmanager", "DB Status", "4. Connected")

	var extra []common.ValidatorInfo
	for _, validator := range podc_global.DBDataList {

		var num uint64

		if podc_global.QRNDDeviceStat == true{
			//log.Info("QRND " ,  " Random Nums" , podc_global.QRNDDeviceStat)
 			randomIndex := rand.Intn(12280) + 1
			log.Info("QRND " ,  " Random Number" , randomIndex)

			num = podc_global.RandomNumbers[randomIndex]


		} else {
			//log.Info("Suedo Random "  ,   " Random Nums", podc_global.QRNDDeviceStat)
			num = rand.Uint64()
		}
		validatorInfo := common.ValidatorInfo{}
		validatorInfo.Address = common.HexToAddress(validator.Address)
		validatorInfo.Qrnd = num
		validatorInfo.Tag = common.Tag(validator.Tag)

		extra = append(extra, validatorInfo)

	}

	return extra

}


func  TestQmanServer (w http.ResponseWriter, req *http.Request){

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Info(string(body))


	m := podc_global.Message{
		Message: "Success",
		Code: http.StatusOK,
	}

	json.NewEncoder(w).Encode(m)



}