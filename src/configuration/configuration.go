package configuration

import (
	"os"
	"user-service/src/servers/manager"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerManagerConfig manager.ServerManagerConfig `json:"serverManager" yaml:"serverManager"`
}

var defaultConfigFile = "config.yaml"

func ReadEnvConfig() *Config {
	// TODO read in config from file mount
	return readDefaultConfig()
}

func readDefaultConfig() *Config {
	config := Config{}
	configBytes, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
