package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type cammieConfig struct {
	BlockedNames     []string `yaml:"blocked_names"`
	BlockedIps       []string `yaml:"blocked_ips"`
	MaxMessageLength int      `yaml:"max_message_length"`
}

type spotifyConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type tapConfig struct {
	URL             string   `yaml:"url"`
	TimestampLayout string   `yaml:"timestamp_layout"`
	Beers           []string `yaml:"beer"`
}

type Config struct {
	Cammie  cammieConfig  `yaml:"cammie"`
	Spotify spotifyConfig `yaml:"spotify"`
	Tap     tapConfig     `yaml:"tap"`
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		configInstance = &Config{}
		data, err := os.ReadFile("config.yaml")
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		if err := yaml.Unmarshal(data, configInstance); err != nil {
			log.Fatalf("Failed to unmarshal config: %v", err)
		}
	})
	return configInstance
}
