package config

// total configurations
type Configurations struct {
	// configuration for your local
	Local	EnvConfig	`json:"local"`
	// configuration for development
	Dev		EnvConfig	`json:"development"`
	// configuration for production
	Prod	EnvConfig	`json:"production"`
}

// configurations of each environment
type EnvConfig struct {
	Consensus	ConsensusConfig	`json:"consensus"`
	Token		TokenConfig		`json:"token"`
	Bootnodes	[]string		`json:"bootnodes"`
	QManagers	[]string		`json:"qmanagers"`
	//Senatornodes []string       `json:"Senatornodes"`
	//Candidatenodes []string     `json:"Candidatenodes"`

}

// detail configurations for consensus part
type ConsensusConfig struct {
	Criteria 	uint		`json:"criteria"`
}

// detail configurations for token economy part
type TokenConfig struct {

}
