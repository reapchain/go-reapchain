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
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/qManager"

	//"github.com/ethereum/go-ethereum/qManager"
)
func main() {
	var (
		listenAddr  = flag.String("addr", ":30301", "listen address")
		genKey      = flag.String("genkey", "", "generate a node key")
		qmanKeyFile = flag.String("qmankey", "", "private key filename")
		qmanKeyHex  = flag.String("qmankeyhex", "", "private key as hex (for testing)")
		qmanKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()

	switch {
	case *genKey != "":
		qmanKey, err = crypto.GenerateKey()  //private key와 public key를 쌍으로 생성하고, private key 구조체는 public key를 담고 있다.

		fmt.Printf("qmanKey(Private key)= %v\n, public Key= %v\n", qmanKey, qmanKey.Public() )
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, qmanKey); err != nil {
			utils.Fatalf("%v", err)
		}
		return
	case *qmanKeyFile == "" && *qmanKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *qmanKeyFile != "" && *qmanKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *qmanKeyFile != "":
		if qmanKey, err = crypto.LoadECDSA(*qmanKeyFile); err != nil {
			fmt.Printf("%v\n", qmanKey )
		}

		//			fmt.Printf("read nodekey= %x\n, read nodekey(Public)= %x\n", nodeKey, nodeKey.Public() )
		if(err != nil) {
			fmt.Printf("%v\n", qmanKey )

			utils.Fatalf("-nodekey: %v", err)
		}

	case *qmanKeyHex != "":
		if qmanKey, err = crypto.HexToECDSA(*qmanKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
		fmt.Printf("return private key from nodekey= %x\n", listenAddr )

	}


		//var account common.Address
		//account = PubkeyToAddress(nodeKey.PublicKey)
		//fmt.Printf("Address(20byte account) : %v\n, %x\n", PubkeyToAddress(nodeKey.PublicKey),account )
		qManager.Start(listenAddr, qmanKey)



	select {}
}
