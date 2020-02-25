package common

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb"
	"net"
)


var (

	QManagerStorage *leveldb.DB
	QManConnected bool
	QManagerAddress *net.UDPAddr
	BootNodeAddress *net.UDPAddr
	QManagerNodeIDStr string
	QManReady bool
	BootNodeReady bool
	BootNodeQManAddr string
)
//

type (
	QManDBStruct struct {
		ID   string
		Address string
	}
)
func ConnectDB() {

	var err error
	QManagerStorage, err = leveldb.OpenFile("level", nil)
	if err != nil{
		log.Info("DB ERROR", "err = ", err)
	}

}

func IsQmanager() (isQMan bool){
	return QManConnected
}

func Find(nodeId []byte) ( found bool) {
	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if QManConnected {
		//var nodeIDString string
		//decode_err := rlp.DecodeBytes(nodeId, &nodeIDString)
		//
		//if decode_err != nil{
		//	log.Info("QManager", "RLP Decode Error = ", decode_err)
		//	return
		//}
		//
		//log.Info("Decoded ID", "qman = ", nodeIDString)

		_, err := QManagerStorage.Get(nodeId, nil)
		if err != nil {
			log.Info("QManager", "DB --", "Node Not Found")

			return false
		}
		//fmt.Println(data)

		//dec := gob.NewDecoder(bytes.NewReader(data))
		//
		//dec.Decode(node)

		//defer QManagerStorage.Close()
		return true

	}

	return true

}


func Save(ID []byte, NodeDetails []byte) ( saved bool) {
	if QManConnected {
		err := QManagerStorage.Put(ID, NodeDetails, nil)

		if err != nil {
			return false
		}
		return true

	}
	return true
}