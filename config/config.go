package config

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"log"
)

type Config struct {
	IncomingPort int
	InfoPort     int
	MetronPort   int
	DebugPort    int
}

func Configuration(configFilePath string) *Config {
	if configFilePath == "" {
		log.Panicln("Missing configuration file path.")
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Panicln("Error loading file:", err)
	}

	config := &Config{}
	if err := json.Unmarshal(configFileBytes, config); err != nil {
		log.Panicln("Error unmarshalling configuration:", err)
	}

	if err = config.validateProperties(); err != nil {
		log.Panicln("Error validating configuration:", err)
	}

	return config
}

func (c *Config) validateProperties() error {
	if c.MetronPort == 0 {
		return errors.New("Metron Port is a required property for the Bosh-HM-Forwarder")
	}

	return nil
}
