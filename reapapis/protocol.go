package reapapis

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

type ProxyCallGuid struct {
	Category string `json:"category"` //iot|did|...
}

type ProxyCall struct {
	Guid    ProxyCallGuid `json:"gid"`
	Service string        `json:"service,omitempty"`
	Call    string        `json:"call"`
	Account string        `json:"account"`
	Value   uint          `json:"value,omitempty"`
	To      string        `json:"to,omitempty"`
}

func (s *ProxyCall) Serialize() []byte {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		os.Exit(-1)
	}
	return jsonBytes
}

func Deserialize(stream string, into *ProxyCall) error {
	err := json.Unmarshal([]byte(stream), &into)
	if err != nil {
		log.Error("failed to Deserialize", "json", err)
		return err
	}
	return nil
}
