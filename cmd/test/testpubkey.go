package main

import (
	"crypto/ecdsa"
	//"encoding/hex"

	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/log"
	"math/big"
	//"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/consensus/istanbul"

	"github.com/ethereum/go-ethereum/discover"
)

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub *ecdsa.PublicKey) discover.NodeID {
	var id discover.NodeID
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
	if nodeKey, err = crypto.HexToECDSA("0d53d73629f75adcecd6fc1eb1c1ecb1e6a20e82a2227c0905b5bc0440be6036"); err != nil {
		fmt.Printf("-nodekeyhex: %v", err)
	}
	fmt.Printf("return private key from nodekey= %x\n", nodeKey )
	fmt.Printf("PublicKey = %v\n", discover.PubkeyID(&nodeKey.PublicKey))

	var account common.Address
	QmanEnode := qmanager[0].ID[0:63]
	// var QmanNodeID discover.NodeID

	var QmanNodeID *discover.NodeID
	QmanNodeID = QmanEnode

	// QmanNodeID = qmanager[0].ID[:]  //배열에서 슬라이스 만들기


	if QmanPublicKey , err  = &discover.QmanNodeID.Pubkey() ; err != nil {
		fmt.Printf("error=%v", err )
	}
	fmt.Printf("PublicKey = %v\n", QmanPublicKey)
	fmt.Printf("Address(20byte account) : %v\n, %x\n", crypto.PubkeyToAddress(nodeKey.PublicKey),account )  //Pulickey -> 20 byte account
	c.qmanager = crypto.PublicKeyBytesToAddress(QmanEnode) //common.Address output from this               //slice ->  "
	fmt.Printf("Address(20byte account)=c.qmanager  : %v\n, %x\n", c.qmanager , c.qmanager )  //Pulickey -> 20 byte account




}
