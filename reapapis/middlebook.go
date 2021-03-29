package reapapis

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"path/filepath"
)

type MiddleBooks struct {
	ks    *keystore.KeyStore
	books map[string]accounts.Account
}

func NewMiddleBooks() *MiddleBooks {
	return &MiddleBooks{}
}

// H/L의 계정을 이용하여 랜덤계정을 생성,
// private key들을 저장할 서버 디렉토리 지정한다.
func (m *MiddleBooks) CreateMiddleBooks(path string) error {
	status, err := IsWritable(path)
	if !status {
		return err
	}
	m.ks = keystore.NewKeyStore(filepath.Join(path, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	return nil
}

// H/L의 계정(id)에 매칭되는 Reapchain 계정리턴 함
func (m *MiddleBooks) Find(id string) (accounts.Account, bool) {
	account, found := m.books[id]
	return account, found
}

// H/L의 계정(id)에 매칭되는 Reapchain 계정을 생성.
// 생성 후, 장부(Books)에 기록함.
func (m *MiddleBooks) Create(id string) (accounts.Account ,error) {
	newAccount, err := m.ks.NewAccount("")
	if err != nil {
		return newAccount, err
	}
	m.books[id] = newAccount
	fmt.Println("ID : ", id, "=new=> : ", newAccount)
	return newAccount, nil
}
