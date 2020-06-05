// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Ethereum Discovery Protocol.
package main

import (
	"crypto/ecdsa"
    "flag"
	"net"
	"strings"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/qManager"

)
// setNodeKey creates a node key from set command line flags, either loading it
// from a file or as a specified hex value. If neither flags were provided, this
// method returns nil and an emphemeral key is to be generated.
/* func setNodeKey(ctx *cli.Context, cfg *p2p.Config) {
	var (
		hex  = ctx.GlobalString(NodeKeyHexFlag.Name)
		file = ctx.GlobalString(NodeKeyFileFlag.Name)
		key  *ecdsa.PrivateKey
		err  error
	)
	switch {
	case file != "" && hex != "":
		Fatalf("Options %q and %q are mutually exclusive", NodeKeyFileFlag.Name, NodeKeyHexFlag.Name)
	case file != "":
		if key, err = crypto.LoadECDSA(file); err != nil {  //산출되는 key 가 20바이트 Address 인가?
			Fatalf("Option %q: %v", NodeKeyFileFlag.Name, err)  //
		}
		cfg.PrivateKey = key
	case hex != "":
		if key, err = crypto.HexToECDSA(hex); err != nil {
			Fatalf("Option %q: %v", NodeKeyHexFlag.Name, err)
		}
		cfg.PrivateKey = key
	}
} */

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d252b0d7c0f7aae9aa624
               c0001a6a20

               5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d56506783588e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c
 * Address: 0x26ea01d982c8aeafab7f440084f90fe078cba92


// nodekey  : 0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036

 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}
/*
func newTestValidatorSet(n int) (istanbul.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
//	vset := validator.NewSet(addrs, istanbul.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}
*/

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
/*
func newSimpleBackend() (backend *simpleBackend, validatorKeys Keys, validatorSet istanbul.ValidatorSet) {
	key, _ := generatePrivateKey()
	validatorSet, validatorKeys = newTestValidatorSet(5)
	backend = &simpleBackend{
		privateKey: key,
		logger:     log.New("backend", "simple"),
		commitCh:   make(chan *types.Block, 1),
	}
	return
}
*/
// FromECDSA exports a private key into a binary dump.
/* func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}
// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
} */
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}
/*
type conn struct {
	fd net.Conn
	p2p.transport
	flags connFlag
	cont  chan error      // The run loop uses cont to signal errors to setupConn.
	id    discover.NodeID // valid after the encryption handshake
	caps  []Cap           // valid after the protocol handshake
	name  string          // valid after the protocol handshake
} */

// Peer represents a connected remote node.
/* type Peer struct {
	rw      *p2p.conn
	running map[string]*p2p.protoRW
	log     log.Logger
	//created mclock.AbsTime

	wg       sync.WaitGroup
	protoErr chan error
	closed   chan struct{}
//	disc     chan DiscReason
} */
/*
func  sendEvent2(event istanbul.ConsensusDataEvent) {
	// FIXME: it's inefficient because it retrieves all peers every time
	//p := pm.findPeer(event.Target)
	var p  *p2p.Peer=&{
		rw : 1
	}

	if p == nil {
		log.Warn("Failed to find peer by address", "addr", event.Target)
		return
	}
	p2p.Send(p.rw, 0, event.Data)
} */
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func main() {
	var (
		listenAddr  = flag.String("addr", ":30391", "listen address")
		port        = flag.String("port", ":30391", "port")
		genKey      = flag.String("genkey", "", "generate a node key")
		writeAddr   = flag.Bool("writeAddress", false, "write out the node's pubkey hash and quit")
		writeAccount   = flag.Bool("writeAccount", false, "write out the node's account hash and quit")
		nodeKeyFile = flag.String("nodekey", "", "private key filename")
		nodeKeyHex  = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc     = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
		netrestrict = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
		runv5       = flag.Bool("v5", false, "run a v5 topic discovery bootnode")
		verbosity   = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		vmodule     = flag.String("vmodule", "", "log verbosity pattern")
		qman        = flag.String( "qman", "", "boot as bootnode and qmanager")

		nodeKey *ecdsa.PrivateKey
		//key  *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	//glogger.Qmanager()
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)

	natm, err := nat.Parse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey()  //private key와 public key를 쌍으로 생성하고, private key 구조체는 public key를 담고 있다.

		fmt.Printf("nodeKey(Private key)= %v\n, public Key= %v\n", nodeKey, nodeKey.Public() )
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			utils.Fatalf("%v", err)
		}
		return
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			fmt.Printf("%v\n", nodeKey )
		}

//			fmt.Printf("read nodekey= %x\n, read nodekey(Public)= %x\n", nodeKey, nodeKey.Public() )
		if(err != nil) {
			fmt.Printf("%v\n", nodeKey )

			utils.Fatalf("-nodekey: %v", err)
		}

	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
		fmt.Printf("return private key from nodekey= %x\n", nodeKey )

	}

	if *writeAddr {
		fmt.Printf("%v\n", discover.PubkeyID(&nodeKey.PublicKey))
		os.Exit(0)
	}
	if *writeAccount {
		var account common.Address
		account = PubkeyToAddress(nodeKey.PublicKey)
		fmt.Printf("Address(20byte account) : %v\n, %x\n", account ,account )
		// sendEvent sends a p2p message with given data to a peer


		//var accountfile string
		//if err = crypto.SaveECDSA(accountfile, account); err != nil {
		//	utils.Fatalf("%v", err)
		//}
		//os.Exit(0)
	}
	var restrictList *netutil.Netlist
	if *netrestrict != "" {
		restrictList, err = netutil.ParseNetlist(*netrestrict)
		if err != nil {
			utils.Fatalf("-netrestrict: %v", err)
		}
	}
	if( *qman !="" ){
		fmt.Printf("%v\n", qman )
		//p2p.Send(p.rw, 0, event.Data)

	}
	if( *listenAddr ==""){
		*listenAddr = GetLocalIP()

		log.Info("I've get Local IP on this machine!", "listenAddr", listenAddr)
	}
	if( *port ==""){
		*port = "30391"

		log.Info("I've get Local IP on this machine!")
	}


	if *runv5 {
		if _, err := discv5.ListenUDP(nodeKey, *listenAddr, natm, "", restrictList); err != nil {
			utils.Fatalf("%v", err)
		}
	} else {  //main Listen server init and run discovery

		var account common.Address
		account = PubkeyToAddress(nodeKey.PublicKey)
		fmt.Printf("Address(20byte account) : %v\n, %x\n", PubkeyToAddress(nodeKey.PublicKey),account )
		qManager.IsBootNode = true
		if _, err := discover.ListenUDP(nodeKey, *listenAddr, natm, "", restrictList); err != nil {  //main
			utils.Fatalf("%v", err)
		}
	}

	select {}
}
