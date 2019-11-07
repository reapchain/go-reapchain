
package main

import (
"encoding/json"
"fmt"
)

type Author struct {
	Name  string
	Email string
}

type Comment struct {
	Id      uint64
	Author  Author // Author 구조체
	Content string
}

type Article struct {
	Id         uint64
	Title      string
	Author     Author    // Author 구조체
	Content    string
	Recommends []string  // 문자열 배열
	Comments   []Comment // Comment 구조체 배열
}

/* const (
	Senator         uint64 = iota // 상원
	Parliamentarian               // 하원
	Candidate                     // 하원, 운영위 후보군
	General                       // 일반 노드, 상원, 하원도 아닌.
	QManager                      // Q-Manager
) */
// ExtraData from Validators Set from Qman
// Validater List from Qmanager


type ValidatorElement struct {
	address []byte
	tag     uint64
	qrnd    uint64

}

//self=enode://095780bcb9ce75c6e97b96daa710262c8eccad48a26636334f4c089f3d08696b6568b9da77b449553cab33eedda29370679a500874b2e66d83bfb268222ebd52@125.131.89.105:30303
//INFO [11-07|14:09:52] RLPx listener up                         self=enode://095780bcb9ce75c6e97b96daa710262c8eccad48a26636334f4c089f3d08696b6568b9da77b449553cab33eedda29370679a500874b2e66d83bfb268222ebd52@125.131.89.105:30303
//INFO [11-07|14:09:52] Mapped network port                      proto=udp extport=30303 intport=30303 interface="UPNP IGDv2-IP1"
//INFO [11-07|14:09:52] IPC endpoint opened: /Users/yongilchoi/Library/Ethereum/geth.ipc
//INFO [11-07|14:09:52] Mapped network port                      proto=tcp extport=30303 intport=30303 interface="UPNP IGDv2-IP1"


// type ValidatorsList []ValidatorElement  //Validator Element 의 구조체 배열,,

func main() {
	/* doc := `
	[{
		"address": "095780bcb9ce75c6e97b96daa710262c8eccad48a26636334f4c089f3d08696b6568b9da77b449553cab33eedda29370679a500874b2e66d83bfb268222ebd52@125.131.89.105:30303"
		"tag": "Senator",
		"qrnd": "1234567890"

	}]
	` */

	doc := `
	[{
		"address": "095780bcb9ce75c6e97b96daa710262c8eccad48a26636334f4c089f3d08696b6568b9da77b449553cab33eedda29370679a500874b2e66d83bfb268222ebd52@125.131.89.105:30303"
		"tag": 1,
		"qrnd": 1234567890
		
	}]
	`


	var data []ValidatorElement // JSON 문서의 데이터를 저장할 구조체 슬라이스 선언


	//json.Unmarshal([]byte(ValidatorList), &data) // doc의 내용을 변환하여 data에 저장
	json.Unmarshal([]byte(doc), &data) // doc의 내용을 변환하여 data에 저장

	fmt.Println(data) // [{1 Hello, world! {Maria maria@exa... (생략)

	//json.Unmarshal([]ValidatorElement(doc), &data) // doc의 내용을 변환하여 data에 저장
	fmt.Println(data) // [{1 Hello, world! {Maria maria@exa... (생략)



}