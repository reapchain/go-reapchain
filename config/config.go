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

	confenv := os.Getenv(env)
	if confenv == "" {
		log.Info("Configurations","REAPCHAIN_ENV environment var not set",  "'local' is Default")
	}else {
		log.Info("Configurations","Current ", confenv)

	}

	var filename string
	if setupInfo := os.Getenv(setupenv); setupInfo =="" {
		log.Info("Configurations Failure", "SETUP_INFO environment var not set", "Please set SETUP_INFO var")
		pwd, err := os.Getwd()
		log.Info("Configurations", "Current Directory: ", pwd)
		if err != nil {
			log.Error("Configurations Failure", "Error Message: ", err)
		}
		filename = pwd + "/config.json"
	} else {
		log.Info("Configurations","SETUP_INFO: ", os.Getenv( "SETUP_INFO"))
		filename = os.Getenv( "SETUP_INFO") + "/config.json"
	}

   // log.Info("confi.json is located at ", "filename", filename )
	log.Info("Configurations","current config file: ",  filename)

	jsonFile, err := os.Open(filename)

	if err != nil {
		log.Error("Configurations Failure","Error Message: ", err)

	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &TotalConfig)

	if confenv == "development" {
		c.Consensus = TotalConfig.Dev.Consensus
		c.Token = TotalConfig.Dev.Token
		c.Bootnodes = TotalConfig.Dev.Bootnodes
		c.QManagers = TotalConfig.Dev.QManagers
		c.Senatornodes = TotalConfig.Dev.Senatornodes
		c.Candidatenodes = TotalConfig.Dev.Candidatenodes

	} else if confenv == "production" {
		c.Consensus = TotalConfig.Prod.Consensus
		c.Token = TotalConfig.Prod.Token
		c.Bootnodes = TotalConfig.Prod.Bootnodes
		c.QManagers = TotalConfig.Prod.QManagers
		c.Senatornodes = TotalConfig.Prod.Senatornodes
		c.Candidatenodes = TotalConfig.Prod.Candidatenodes

	} else {
		c.Consensus = TotalConfig.Local.Consensus
		c.Token = TotalConfig.Local.Token
		c.Bootnodes = TotalConfig.Local.Bootnodes
		c.QManagers = TotalConfig.Local.QManagers
		c.Senatornodes = TotalConfig.Local.Senatornodes
		c.Candidatenodes = TotalConfig.Local.Candidatenodes

	}
	//fmt.Printf("\nSenator   nodes : %s\n", c.Senatornodes)
	//fmt.Printf("\nCandidate nodes : %s\n", c.Candidatenodes)


}
