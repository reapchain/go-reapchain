package global

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
	"net"
)

var (

	QManagerStorage *leveldb.DB
	QManConnected bool  //If I'm Qmanager, this is true. when start p2p node, check whether Qman or not.
	QManagerAddress *net.UDPAddr
	BootNodeReady bool
	QRNGDeviceStat bool
	QRNGFilePrefix string
	RandomNumbers []uint64
	BootNodePort int
	BootNodeID string
	IsBootNode bool
	QManPubKey *ecdsa.PublicKey
	GovernanceList []GovStruct
	DBDataList []QManDBStruct
)
type (
	QManDBStruct struct {
		ID      string
		Address  string
		Timestamp	string
		Tag string
	}
)
type Message struct {
	Message string
	Code int
}
type GovStruct struct {
	Validator string
	Tag common.Tag
}
type RequestStruct struct {
	Proposer string
}
type RequestCoordiStruct struct {
	QRND uint64
}
type CoordiDecideStruct struct {
	Status bool
}
func CheckBootNodePortAndID(NodeID string, Port int) bool{
	 if NodeID == BootNodeID && Port == BootNodePort {
	 	return true
	 }
	 return false
}