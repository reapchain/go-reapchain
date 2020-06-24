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
		fmt.Printf("REAPCHAIN_ENV environment var not set, local will set\n")
	}
	pwd, err := os.Getwd()
	fmt.Printf("current working directory: pwd= %v \n",  pwd)
	if err != nil {
		fmt.Printf("failed to get current working directory: pwd= %v , err=%v",  pwd, err)
	}
    //var filename string
    filename := pwd + "/setup_info/config.json"
   // log.Info("confi.json is located at ", "filename", filename )
	fmt.Printf("\ncurrent config file : %v ",  filename)

	jsonFile, err := os.Open(filename)

	if err != nil {
		log.Error("Failed to read configurations", "error", err)
		panic(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &TotalConfig)

	if env == "development" {
		c.Consensus = TotalConfig.Dev.Consensus
		c.Token = TotalConfig.Dev.Token
		c.Bootnodes = TotalConfig.Dev.Bootnodes
		c.Senatornodes = TotalConfig.Local.Senatornodes
	} else if env == "production" {
		c.Consensus = TotalConfig.Prod.Consensus
		c.Token = TotalConfig.Prod.Token
		c.Bootnodes = TotalConfig.Prod.Bootnodes
		c.Senatornodes = TotalConfig.Local.Senatornodes
	} else {
		c.Consensus = TotalConfig.Local.Consensus
		c.Token = TotalConfig.Local.Token
		c.Bootnodes = TotalConfig.Local.Bootnodes
		c.Senatornodes = TotalConfig.Local.Senatornodes
	}
	fmt.Printf("\nSenator nodes : %s\n", c.Senatornodes)


}
