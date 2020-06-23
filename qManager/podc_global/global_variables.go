package podc_global

import (
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
)

type (
	QManDBStruct struct {
		ID      string
		Address  string
		Timestamp	string
	}
)