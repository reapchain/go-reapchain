package qManager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"log"
	"net/http"
	//"encoding/hex"
	"io/ioutil"
	// "bufio"
	// "crypto/x509"
	// "os"
	// "encoding/pem"
	// "crypto/ecdsa"
	//  "bytes"

	//"crypto/rand"
	"fmt"
	//"github.com/ethereum/go-ethereum/crypto/ecies"
	"time"
)



type Message struct {
	Message string
	Code int
}

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
	log.Println(body)
	log.Println(string(body))



	var govStruct []podc_global.GovStruct
	err = json.Unmarshal(body, &govStruct)
	if err != nil {
		m := Message{
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

	m := Message{
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

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp)

	secondTime := time.Now().Add(time.Second * time.Duration(12)).Format("2006-01-02 15:04:05")

	fmt.Println(secondTime)


	t, _ := time.Parse("2006-01-02 15:04:05", timestamp )
	t2, _ := time.Parse("2006-01-02 15:04:05", secondTime )
	fmt.Println(t)



	diff := t2.Sub(t)
	fmt.Println(diff)
	// s := string(bodyBytes[:])
	// log.Println(s)
	// myString := hex.EncodeToString(bodyBytes)
	// log.Println(myString)
	//http.HandleFunc("/hello", hello)
	http.HandleFunc("/RequestQmanager", RequestQmanager)
	//addr, err := net.ResolveUDPAddr("udp", *Addr)

	http.ListenAndServe(*Addr, nil)
}


func MakeRequest() {

	message := map[string]interface{}{
		"hello": "world",
		"life":  42,
		"embedded": map[string]string{
			"yes": "of course!",
		},
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("http://192.168.0.67:5050/RequestQmanager", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
}