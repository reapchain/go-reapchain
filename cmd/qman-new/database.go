package main

import (
	"bytes"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

var (
	lastSeqPrefix   = []byte("lastSeq-")
	lastSetPrefix   = []byte("lastSet-")
	senatorPrefix   = []byte("senator-")
	candidatePrefix = []byte("candidate-")
)

func ReadSenators(db ethdb.Database) []common.Address {
	prefix := senatorPrefix
	keylen := len(prefix) + len(common.Address{})
	it := db.(*ethdb.LDBDatabase).NewIterator()
	defer it.Release()

	var result []common.Address

	it.Seek(prefix)
	for bytes.HasPrefix(it.Key(), prefix) {
		var addr common.Address
		copy(addr[:], it.Key()[keylen-len(addr):])
		//value := it.Value()
		result = append(result, addr)
		it.Next()
	}
	return result
}

func ReadCandidates(db ethdb.Database) []common.Address {
	prefix := candidatePrefix
	keylen := len(prefix) + len(common.Address{})
	it := db.(*ethdb.LDBDatabase).NewIterator()
	// it := db.NewIterator(prefix, nil)
	defer it.Release()

	var result []common.Address

	it.Seek(prefix)
	for bytes.HasPrefix(it.Key(), prefix) {
		var addr common.Address
		copy(addr[:], it.Key()[keylen-len(addr):])
		//value := it.Value()
		result = append(result, addr)
		it.Next()
	}
	return result
}

func WriteSenator(db ethdb.Database, addr common.Address, tag int) error {
	key := append(senatorPrefix, addr[:]...)
	err := db.Put(key, []byte(strconv.Itoa(tag)))
	if err != nil {
		log.Error("db write error", "err", err)
		return err
	}
	return nil
}

func WriteCandidate(db ethdb.Database, addr common.Address, tag int) error {
	key := append(candidatePrefix, addr[:]...)
	err := db.Put(key, []byte(strconv.Itoa(tag)))
	if err != nil {
		log.Error("db write error", "err", err)
		return err
	}
	return nil
}

func RemoveSenator(db ethdb.Database, addr common.Address) error {
	key := append(senatorPrefix, addr[:]...)
	err := db.Delete(key)
	if err != nil {
		log.Error("db delete error", "err", err)
		return err
	}
	return nil
}

func RemoveCandidate(db ethdb.Database, addr common.Address) error {
	key := append(candidatePrefix, addr[:]...)
	err := db.Delete(key)
	if err != nil {
		log.Error("db delete error", "err", err)
		return err
	}
	return nil
}
