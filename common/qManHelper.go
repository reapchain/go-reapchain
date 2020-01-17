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
//type (
//	QManDBStruct struct {
//		ID   discover.NodeID
//		pubKey *ecdsa.PublicKey
//		address common.Address
//	}
//)
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

func Find(nodeId string) ( found bool) {
	//QManagerStorage, err = leveldb.OpenFile("level", nil)

	//var data []byte
	if QManConnected {
		_, err := QManagerStorage.Get([]byte(nodeId), nil)
		log.Info("Node Not Found", "err = ", err)
		if err != nil {
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


func Save(ID []byte, address []byte) ( saved bool) {
	if QManConnected {
		err := QManagerStorage.Put(ID, address, nil)

		if err != nil {
			return false
		}
		return true

	}
	return true
}