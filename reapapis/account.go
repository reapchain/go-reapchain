package reapapis

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"path/filepath"
)

// 테스트용 Account Address.
const GovernanceAddress = "0xe10f734218c24be1d64864c25f3ae2f1d45ceda8"

type Account struct {
	ks      *keystore.KeyStore
	account accounts.Account
}

func NewImportAccount() *Account {
	return &Account{}
}

func (g *Account) Account() accounts.Account {
	return g.account
}

// Account private key가 keystore import 한다.
// ex) /node/keystore/UTC--2020-11..... 일 경우, Import("/node")
func (g *Account) Import(path string) error {
	keystorePath := path + "/" + "keystore"
	status, err := IsWritable(keystorePath)
	if !status {
		return err
	}
	g.ks = keystore.NewKeyStore(filepath.Join(path, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)

	if len(g.ks.Accounts()) > 1 {
		return errors.New("too many keys in keystore designed")
	}
	g.account = g.ks.Accounts()[0]

	log.Debug("Governance Account Info", "Governance Address", GovernanceAddress)
	log.Debug("Import Account Info", "Account Address",g.account.Address.String())

	return nil
}

// Account 계정으로 Tx를 sign 한다.
func (g *Account) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return g.ks.SignTx(g.account, tx, chainID)
}

// Account 계정으로 Tx를 sign한다. 패스워드 값필요.
func (g *Account) SignTxWithPassphrase(tx *types.Transaction, passphrase string, chainID *big.Int) (*types.Transaction, error) {
	return g.ks.SignTxWithPassphrase(g.account, passphrase, tx, chainID)
}
