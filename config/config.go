package config

import (
	//"errors"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/ethereum/go-ethereum/log"

)

var Config EnvConfig
var TotalConfig Configurations

func (c *EnvConfig) GetConfig(env string) {
	jsonFile, err := os.Open("config.json")
	defer jsonFile.Close()

	if err != nil {
		log.Error("Failed to read configurations", "error", err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &TotalConfig)

	if env == "development" {
		c.Consensus = TotalConfig.Dev.Consensus
		c.Token = TotalConfig.Dev.Token
		c.Bootnodes = TotalConfig.Dev.Bootnodes
	} else if env == "production" {
		c.Consensus = TotalConfig.Prod.Consensus
		c.Token = TotalConfig.Prod.Token
		c.Bootnodes = TotalConfig.Prod.Bootnodes
	} else {
		c.Consensus = TotalConfig.Local.Consensus
		c.Token = TotalConfig.Local.Token
		c.Bootnodes = TotalConfig.Local.Bootnodes
	}
}
