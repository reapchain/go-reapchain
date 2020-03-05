package config

import (
	//"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/ethereum/go-ethereum/log"

)

var Config EnvConfig
var TotalConfig Configurations

func (c *EnvConfig) GetConfig(env string) {

	if confenv := os.Getenv(env); confenv == "" {
		fmt.Printf("REAPCHAIN_ENV environment var not set, local will set")
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current working directory: pwd= %v , err=%v",  pwd, err)
	}


	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Error("Failed to read configurations", "error", err)
	}
	log.Info("path name of config.json =", "jsonFile", jsonFile ) //added by yichoi to check directory path of config.json

	defer jsonFile.Close()



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
