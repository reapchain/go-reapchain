// Raapchain paper

/*
1. Frontnode role : 프론트 노드는 블록 생성을 요청하는 새로운 Tx를 수신하고 Bx를 생성하는 시점에 Qmanager 서버세 접속해서 퀀텀난수 생성을 요청
   (태헌)

2. Qmanager role : Qmanager는 프론트 노드의 요청을 받고 블록생성 작업과 합의를 위한 코디 및 운영위 후보군과 이들의 인식을 위한 퀀텀난수 정보를
                   암호화하여 이를 프론트 노드에 전송. 선정된 코디에게는 운영위 후보군 정보를 전송
   (마틴)

3. 프론트 노드는 Qmanager로부터 전송받은 암호화 정보로 블록 Header에 Extra data로 구성하여 전체 노드에 브로드캐스팅
   (태헌)

4. 각 노드는 전송받은 블록의 Extra data에서 자기자신의 해시를 검증하고 자신의 정보인 경우에만 복호화하여 코디와 후보군 확인
   (태헌)

*/

package quantum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"io/ioutil"
	"os"
	"time"
	//"github.com/ethereum/go-ethereum/crypto"
	//"errors"
)

const bufferSize = 16384
//const filePrefix = "/run/media/user/E8EE-1C60/"
const filePrefix ="/Volumes/PSoC USB/"  //on my MAC

//func GenerateQrnd(nodeid []*discover.NodeID) []byte {

// 백서에 노드ID 주면 양자난수를 리턴하는 것으로 명시되어 있음.
// 노드ID 상관 없이 양자난수 발생기 만들어야 여기서
func GenerateQrnd(nodeid  *discover.NodeID) []byte  {
	var ioldIndex byte
	var Qrnddata []byte

	//QmanEnode := qmanager[0].ID[:]  //여기까지 정상
	//nodeid :=

	//c.qmanager = crypto.PublicKeyBytesToAddress(QmanEnode) //common.Address output from this [account addr]              //slice ->

	//Qmanager account address(20byte): 926ea01d982c8aeafab7f440084f90fe078cba92

	for {
		buffer := readUpFile()   //buffer data type is slice []byte
		if buffer[0] == ioldIndex {
			writeDnFile(buffer[0])
			time.Sleep(1 * time.Second)

		} else {
			ioldIndex = buffer[0]

			if len(buffer) > 0 {
				for j := 1; j <= bufferSize-1; j++ {
					fmt.Printf("[0x%02x]", buffer[j])
				}
				fmt.Println()

				writeDnFile(buffer[0])

				Qrnddata = buffer[1:] //?  //slice to slice data exchange
			} else {
				time.Sleep(1 * time.Second)
			}

			fmt.Println(fmt.Sprintf("INDEX: %d", ioldIndex))
		}
	}


	return Qrnddata
}

func readUpFile() []byte {
	buffer, err := ioutil.ReadFile(filePrefix + "up.ini")
	if err != nil {
		fmt.Println("OPEN-ERROR: up.ini")
		panic(err.Error())
	}

	return buffer
}

func writeDnFile(data byte) {
	os.Remove(filePrefix + "dn.ini")
	err := ioutil.WriteFile(filePrefix+"dn.ini", []byte{data}, 0644)
	if err != nil {
		fmt.Println("CREATE-ERROR: dn.ini")
		panic(err.Error())
	}
}

