package podc_global

import (
	"crypto/ecdsa"
	"github.com/syndtr/goleveldb/leveldb"
	"net"
)

var (

	QManagerStorage *leveldb.DB
	QManConnected bool  //If I'm Qmanager, this is true. when start p2p node, check whether Qman or not.
	QManagerAddress *net.UDPAddr
	BootNodeReady bool
	QRNDDeviceStat bool
	RandomNumbers []uint64
	BootNodePort int
	BootNodeID string
	IsBootNode bool
	QManPubKey *ecdsa.PublicKey
	GovernanceList []GovStruct
)

type (
	QManDBStruct struct {
		ID      string
		Address  string
		Timestamp	string
		Tag uint64
	}
)

type GovStruct struct {
	Validator string
	Tag string
}

func CheckBootNodePortAndID(NodeID string, Port int) bool{
	//log.Info("BootNode", "ID = ", NodeID)
	//log.Info("BootNode", "IP Addr = ", Port)
	//
	//log.Info("Global", "ID = ", BootNodeID)
	//log.Info("Global", "IP Addr = ", BootNodePort)

	 if NodeID == BootNodeID && Port == BootNodePort{
	 	return true
	 }

	 return false
}