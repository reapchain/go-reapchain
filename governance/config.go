// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// yhheo begin
package governance

import (
    "fmt"
	"io"
    "os"
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	// Path within the datadir to the governance's public key
	datadirGovernanceKey = "governancekey"
)

type GovernanceConfig struct {
	// This field must be set to a public key.
	PublicKey  []byte

	Name string

	DataDir string

	Governance bool
}

//var gc = &GovernanceConfig {
//	Name:	datadirGovernanceKey,
var gc = &GovernanceConfig{}

func GetFileName() string {
	gc.Name = datadirGovernanceKey
	return gc.Name
}

func LoadKey(DataDir string, flag bool) {

	if DataDir != "" && flag == true {
		gc.Governance = flag
		gc.DataDir = DataDir

		governanceKey(gc)

		log.Info("Initialised governance configuration", "gconfig", gc)
	}
}

func CheckPublicKey(pubkey []byte) bool {

    var isGovernance bool

	fmt.Printf("\nfunc CheckPublicKey\n gc.Governance = %t\n", gc.Governance)

    if gc.Governance {
        pkHash := common.BytesToHash(pubkey)
        fmt.Printf(" Tx.pubkey    = %x\n", pubkey)
        fmt.Printf(" Tx.pkHash    = %x\n", pkHash)
        fmt.Printf(" gc.PublicKey = %x\n", gc.PublicKey)

        governanceKey(gc)

        isGovernance = bytes.Equal(gc.PublicKey, []byte(pkHash[:]))
    } else {
        isGovernance = false
    }

    fmt.Println(" isGovernance =", isGovernance)

    return isGovernance
}

func governanceKey(gc *GovernanceConfig) []byte {

    // Use any specifically configured key.
    if gc.PublicKey != nil {
        return gc.PublicKey
    }

    if key, err := loadGKey(gc.DataDir); err == nil {
        gc.PublicKey = key
        return key
    } else {
        log.Warn(fmt.Sprintf("Failed to load public key: %v", err))
    }
    return nil
}

// loadGKey loads a public key from the given file.
func loadGKey(file string) ([]byte, error) {
    buf := make([]byte, 64)
    fd, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer fd.Close()
    if _, err := io.ReadFull(fd, buf); err != nil {
        return nil, err
    }

    key, err := hex.DecodeString(string(buf))
    if err != nil {
        return nil, err
    }
    return key, nil
}
// yhheo - end
