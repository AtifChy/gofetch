// Package config
package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//go:embed default_config.json
var defaultConfigFile []byte

type Config struct {
	Logo LogoConfig `json:"logo"`
}

type LogoConfig struct {
	Colors map[int]string `json:"colors"`
}

func LoadDefaultConfig() Config {
	var config Config
	if err := json.Unmarshal(defaultConfigFile, &config); err != nil {
		log.Fatalf("could not parse embed default config: %v", err)
	}
	return config
}

func LoadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("could not read config file %s: %w", path, err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not parse config file %s: %w", path, err)
	}

	return config, nil
}
