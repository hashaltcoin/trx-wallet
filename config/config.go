package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"sync"
)

var (
	config     *Config
	mux        sync.Mutex
	configPath string
)

type Config struct {
	NodeURL   string     `json:"node_url"`
	Contracts []Contract `json:"contracts"`
}

type Contract struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	ContractType string `json:"contract_type"`
}

func init() {
	flag.StringVar(&configPath, "config", "./config.json", "Path to the primary JSON config file")
	flag.Parse()
}

//ParseConfig and set a shared config entry
func ParseConfig(path string) (*Config, error) {
	if len(path) == 0 {
		path = configPath
		if len(path) == 0 {
			panic("Invalid config path. Not provided and not a command line option")
		}
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("Invalid ConfigPath setting: %s", path)
	}
	if info.IsDir() {
		log.Fatalf("ConfigPath is a directory: %s", path)
	}

	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(configFile)
	err = dec.Decode(&config)
	return config, nil
}

//GetConfig returns a shared instance of config
func GetConfig() (*Config, error) {
	if config == nil {
		mux.Lock()
		defer mux.Unlock()
		if config == nil {
			_, err := ParseConfig("")
			if err != nil {
				return nil, err
			}
		}
	}
	return config, nil
}
