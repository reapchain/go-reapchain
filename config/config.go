package config

import (
	//"errors"
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"os"
)

var Config EnvConfig
var TotalConfig Configurations

func (c *EnvConfig) GetConfig(env string, setupenv string) {

	if confenv := os.Getenv(env); confenv == "" {
		log.Info("Configurations","REAPCHAIN_ENV environment var not set",  "local will set")
	}
	if setupInfo := os.Getenv(setupenv); setupInfo ==""{
		log.Info("Configurations","SETUP_INFO environment var not set",  "you should set SETUP_INFO var")
	}
	pwd, err := os.Getwd()
	log.Info("Configurations","Current Directory: ",  pwd)
	if err != nil {
		log.Info("Configurations","failed to get current working directory: ", err)
	}

	log.Info("Configurations","SETUP_INFO: ", os.Getenv( "SETUP_INFO"))
    filename := os.Getenv( "SETUP_INFO") + "config.json"
   // log.Info("confi.json is located at ", "filename", filename )
	log.Info("Configurations","current config file: ",  filename)

	jsonFile, err := os.Open(filename)

	if err != nil {
		log.Error("Configurations","Failed to read configurations", "error", err)

	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &TotalConfig)

	if env == "development" {
		c.Consensus = TotalConfig.Dev.Consensus
		c.Token = TotalConfig.Dev.Token
		c.Bootnodes = TotalConfig.Dev.Bootnodes
		c.QManagers = TotalConfig.Dev.QManagers
		//c.Senatornodes = TotalConfig.Local.Senatornodes
		//c.Candidatenodes = TotalConfig.Local.Candidatenodes

	} else if env == "production" {
		c.Consensus = TotalConfig.Prod.Consensus
		c.Token = TotalConfig.Prod.Token
		c.Bootnodes = TotalConfig.Prod.Bootnodes
		c.QManagers = TotalConfig.Dev.QManagers

		//c.Senatornodes = TotalConfig.Local.Senatornodes
		//c.Candidatenodes = TotalConfig.Local.Candidatenodes
	} else {
		c.Consensus = TotalConfig.Local.Consensus
		c.Token = TotalConfig.Local.Token
		c.Bootnodes = TotalConfig.Local.Bootnodes
		c.QManagers = TotalConfig.Dev.QManagers

		//c.Senatornodes = TotalConfig.Local.Senatornodes
		//c.Candidatenodes = TotalConfig.Local.Candidatenodes
	}
	//fmt.Printf("\nSenator   nodes : %s\n", c.Senatornodes)
	//fmt.Printf("\nCandidate nodes : %s\n", c.Candidatenodes)


}
