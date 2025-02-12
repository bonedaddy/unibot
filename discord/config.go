package discord

import (
	"io/ioutil"
	"os"

	"github.com/bonedaddy/unibot/bclient"
	"gopkg.in/yaml.v2"
)

// Config bundles together discord configuration information
type Config struct {
	// if nil we dont use infura and connect directly to the rpc node below
	InfuraAPIKey    string    `yaml:"infura_api_key"`
	InfuraWSEnabled bool      `yaml:"infura_ws_enabled"`
	ETHRPCEndpoint  string    `yaml:"eth_rpc_endpoint"`
	Watchers        []Watcher `yaml:"watchers"`
	Database        Database  `yaml:"database"`
}

// Database provides configuration over our database connection
type Database struct {
	Type           string `yaml:"type"` // sqlite or postgres, if sqlite all other options except DBName are ignored
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	DBName         string `yaml:"db_name"`
	DBPath         string `yaml:"db_path"`
	SSLModeDisable bool   `yaml:"ssl_mode_disable"`
}

// Watcher is used to start a process that watches the price of a token
// and posts its value as a name
type Watcher struct {
	DiscordToken  string `yaml:"discord_token"`
	Token0Address string `yaml:"token0_address"`
	Token1Address string `yaml:"token1_address"`
	Pair          string `yaml:"pair"`
	Decimals      int    `yaml:"decimals"`
}

var (
	// ExampleConfig is primarily used to provide a template for generating the config file
	ExampleConfig = &Config{
		InfuraAPIKey:    "INFURA-KEY",
		InfuraWSEnabled: false,
		ETHRPCEndpoint:  "http://localhost:8545",
		Watchers: []Watcher{
			{DiscordToken: "CHANGEME-TOKEN", Token0Address: bclient.WETHTokenAddress.String(), Token1Address: bclient.DAITokenAddress.String()},
		},
		Database: Database{
			Type:           "sqlite",
			Host:           "localhost",
			Port:           "5432",
			User:           "user",
			Pass:           "pass",
			DBName:         "indexed",
			DBPath:         "/changeme",
			SSLModeDisable: false,
		},
	}
)

// NewConfig generates a new config and stores at path
func NewConfig(path string) error {
	data, err := yaml.Marshal(ExampleConfig)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}

// LoadConfig loads the configuration
func LoadConfig(path string) (*Config, error) {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(r, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
