package reapapis

import (
	"io/ioutil"
	"os"
	"testing"
)

var yamlText = `
proxy:
  local:
    ratio: 1
    governance: "/home/jaeseong/governance"
    kafka:
      address: "127.0.0.1:9092"
      topic: "reapchain"
    node:
      rpcaddress: "http://192.168.11.1:8541"

scanner:
  local:
    account: "reap"
    passwd: "45b8"
    ip: "test.reappay.net"
    port: 3306
    dbname: "reapdb"
`

const tempFilepath = "/tmp/reapchain.yml"

func TestNewConfigLoad(t *testing.T) {
	err := ioutil.WriteFile(tempFilepath, []byte(yamlText), 0644)
	if err != nil {
		t.Fatal("cant write a test file")
	}
	defer os.Remove(tempFilepath)

	config, err := LoadConfigFile(tempFilepath)
	if err != nil {
		t.Error("LoadConfigFile")
	}

	if config.Proxy.Local.Kafka.Topic != "reapchain" {
		t.Error("not equal")
	}
	if config.Proxy.Local.Kafka.Address != "127.0.0.1:9092" {
		t.Error("not equal")
	}
	if config.Proxy.Local.Ratio != 1 {
		t.Error("not equal")
	}
	if config.Proxy.Local.Governance != "/home/jaeseong/governance" {
		t.Error("not equal")
	}
	if config.Scanner.Local.Ip != "test.reappay.net" {
		t.Error("not equal")
	}
	if config.Scanner.Local.Port != "3306" {
		t.Error("not equal")
	}
	if config.Scanner.Local.Account!= "reap"{
		t.Error("not equal")
	}
	if config.Scanner.Local.Passwd!= "45b8"{
		t.Error("not equal")
	}
	if config.Scanner.Local.DBname!= "reapdb"{
		t.Error("not equal")
	}
}
