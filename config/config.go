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

func (c *EnvConfig) GetConfig(env string, setupenv string) {

	if confenv := os.Getenv(env); confenv == "" {
		fmt.Printf("REAPCHAIN_ENV environment var not set, local will set\n")
	}
	if setupInfo := os.Getenv(setupenv); setupInfo ==""{
		fmt.Printf("SETUP_INFO environment var not set, you should setup_info environment \n")
	}
	pwd, err := os.Getwd()
	fmt.Printf("current working directory: pwd= %v \n",  pwd)
	if err != nil {
		fmt.Printf("failed to get current working directory: pwd= %v , err=%v",  pwd, err)
	}

    fmt.Println("SETUP_INFO :", os.Getenv( "SETUP_INFO"))
    filename := os.Getenv( "SETUP_INFO") + "/config.json"
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
		c.Candidatenodes = TotalConfig.Local.Candidatenodes

	} else if env == "production" {
		c.Consensus = TotalConfig.Prod.Consensus
		c.Token = TotalConfig.Prod.Token
		c.Bootnodes = TotalConfig.Prod.Bootnodes
		c.Senatornodes = TotalConfig.Local.Senatornodes
		c.Candidatenodes = TotalConfig.Local.Candidatenodes
	} else {
		c.Consensus = TotalConfig.Local.Consensus
		c.Token = TotalConfig.Local.Token
		c.Bootnodes = TotalConfig.Local.Bootnodes
		c.Senatornodes = TotalConfig.Local.Senatornodes
		c.Candidatenodes = TotalConfig.Local.Candidatenodes
	}
	fmt.Printf("\nSenator   nodes : %s\n", c.Senatornodes)
	fmt.Printf("\nCandidate nodes : %s\n", c.Candidatenodes)


}
