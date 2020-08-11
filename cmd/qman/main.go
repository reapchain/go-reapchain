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

package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/config"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/qManager"
	"github.com/ethereum/go-ethereum/qManager/podc_global"
	"os"

	//"github.com/ethereum/go-ethereum/qManager"
)
func main() {
	var (
		listenAddr  = flag.String("addr", ":30301", "listen address")
		genKey      = flag.String("genkey", "", "generate a qman key")
		qmanKeyFile = flag.String("qmankey", "", "private key filename")
		qmanKeyHex  = flag.String("qmankeyhex", "", "private key as hex (for testing)")
		verbosity   = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		vmodule     = flag.String("vmodule", "", "log verbosity pattern")

		qmanKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	//glogger.Qmanager()
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)


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
		utils.Fatalf("Use -qmankey or -qmankeyhex to specify a private key")
	case *qmanKeyFile != "" && *qmanKeyHex != "":
		utils.Fatalf("Options -qmankey and -qmankeyhex are mutually exclusive")
	case *qmanKeyFile != "":
		if qmanKey, err = crypto.LoadECDSA(*qmanKeyFile); err != nil {
			fmt.Printf("%v\n", qmanKey )
		}

		if(err != nil) {
			fmt.Printf("%v\n", qmanKey )

			utils.Fatalf("-qmankey: %v", err)
		}

	case *qmanKeyHex != "":
		if qmanKey, err = crypto.HexToECDSA(*qmanKeyHex); err != nil {
			utils.Fatalf("-qmankeyhex: %v", err)
		}
		fmt.Printf("return private key from qmankey= %x\n", listenAddr )

	}

	//pwd, err := os.Getwd()
	//fmt.Printf("current working directory: pwd= %v \n",  pwd)
	//if err != nil {
	//	fmt.Printf("failed to get current working directory: pwd= %v , err=%v",  pwd, err)
	//}

	log.Info("QManager Standalone Started")

	podc_global.QManConnected = true
	config.Config.GetConfig("REAPCHAIN_ENV", "SETUP_INFO")
		//var account common.Address
		//account = PubkeyToAddress(nodeKey.PublicKey)
		//fmt.Printf("Address(20byte account) : %v\n, %x\n", PubkeyToAddress(nodeKey.PublicKey),account )

		qManager.InitializeQManager()
		qManager.Start(listenAddr, qmanKey)





	select {}
}
