package reapapis

import (
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DB struct {
	Ip      string `yaml:"ip"`
	Port    string `yaml:"port"`
	Account string `yaml:"account"`
	Passwd  string `yaml:"passwd"`
	DBname  string `yaml:"dbname"`
}

type Scanner struct {
	Local      DB `yaml:"local,omitempty"`
	Test       DB `yaml:"test,omitempty"`
	Production DB `yaml:"production,omitempty"`
}

type ProxyInfo struct {
	Ratio float32 `yaml:"ratio"`
	Governance string `yaml:"governance,omitempty"`
	Kafka struct {
		Address string `yaml:"address"`
		Topic   string `yaml:"topic"`
	} `yaml:"kafka"`
	Node struct {
		RpcAddress string `yaml:"rpcaddress"`
	} `yaml:"node"`
}

type Proxy struct {
	Local      ProxyInfo `yaml:"local,omitempty"`
	Test       ProxyInfo `yaml:"test,omitempty"`
	Production ProxyInfo `yaml:"production,omitempty"`
}

type Config struct {
	Scanner `yaml:"scanner"`
	Proxy   `yaml:"proxy"`
}

func LoadConfigFile(ABSPath string) (Config, error) {
	var config Config
	filename, err := os.Open(ABSPath)
	yamlFile, err := ioutil.ReadAll(filename)
	if err != nil {
		log.Error("Cant load Config file", ABSPath, err)
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Error("yaml Syntax error", ABSPath, err)
		return config, err
	}
	return config, nil
}
