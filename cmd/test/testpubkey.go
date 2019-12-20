package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"

	//"encoding/hex"

	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/log"
	"math/big"
	//"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/consensus/istanbul"

	"github.com/ethereum/go-ethereum/p2p/discover"
)
type NodeID [512 / 8]byte  // 512/8 = 64byte
// PubkeyID returns a marshaled representation of the given public key.

func PubkeyID(pub *ecdsa.PublicKey) NodeID {
	var id NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

// Pubkey returns the public key represented by the node ID.
// It returns an error if the ID is not a point on the curve.
func (id NodeID) Pubkey() (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: crypto.S256(), X: new(big.Int), Y: new(big.Int)}
	half := len(id) / 2
	p.X.SetBytes(id[:half])
	p.Y.SetBytes(id[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("id is invalid secp256k1 curve point")
	}
	return p, nil
}

func main() {
	// nodekey  : 0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036
	var nodeKey *ecdsa.PrivateKey
	var QmanPublicKey *ecdsa.PublicKey
	var err error
	if nodeKey, err = crypto.HexToECDSA("0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036"); err != nil {
		fmt.Printf("-nodekeyhex: %v\n", err)
	}
	fmt.Printf("Private key            = %x\n", nodeKey )
	fmt.Printf("PublicKey(NodeID type) = %v\n", discover.PubkeyID(&nodeKey.PublicKey))
	//Qmanager=[enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c
	var s string ="5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c"

	var qmanager [64]byte
    var i int8 =0

    for  i := 1; i < 64 ; i++  {
		fmt.Printf("%x", qmanager[i] )
	}

	var account common.Address
	//QmanEnode := qmanager[0].ID[0:63]
	// var QmanNodeID discover.NodeID

	//var QmanNodeID *discover.NodeID
	//QmanNodeID = QmanEnode

	// QmanNodeID = qmanager[0].ID[:]  //배열에서 슬라이스 만들기


	//if QmanPublicKey , err  = &discover.QmanNodeID.Pubkey() ; err != nil {
		fmt.Printf("error=%v", err )
	//}
	fmt.Printf("PublicKey = %v\n", QmanPublicKey)
	fmt.Printf("Address(20byte account) : %v\n, %x\n", crypto.PubkeyToAddress(nodeKey.PublicKey),account )  //Pulickey -> 20 byte account
	//c.qmanager = crypto.PublicKeyBytesToAddress(QmanEnode) //common.Address output from this               //slice ->  "
	//fmt.Printf("Address(20byte account)=c.qmanager  : %v\n, %x\n", c.qmanager , c.qmanager )  //Pulickey -> 20 byte account




}
