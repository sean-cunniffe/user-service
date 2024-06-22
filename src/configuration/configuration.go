package configuration

import (
	"os"
	healthprobes "user-service/src/probes"
	"user-service/src/servers/manager"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config holds configs for instances
type Config struct {
	LogLevel            string                      `json:"logLevel" yaml:"logLevel"` // debug, info, warn, error, fatal, panic
	ServerManagerConfig manager.ServerManagerConfig `json:"serverManager" yaml:"serverManager"`
	ProbeConfig         healthprobes.ProbeConfig    `json:"probeConfig" yaml:"probeConfig"`
}

var defaultConfigFile = "config.yaml"

// ReadEnvConfig reads default config and then overrides values with environment variables
func ReadEnvConfig() *Config {
	// TODO read in config from file mount
	config := readDefaultConfig()
	setLogLevel(config.LogLevel)
	return config
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	}
}

func readDefaultConfig() *Config {
	defaultConfig, exists := os.LookupEnv("DEFAULT_CONFIG")
	if exists {
		defaultConfigFile = defaultConfig
	}

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
