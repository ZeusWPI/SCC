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

type buzzerConfig struct {
	Pin  int    `yaml:"pin"`
	Song string `yaml:"song"`
}

type spotifyConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type tapConfig struct {
	URL   string   `yaml:"url"`
	Beers []string `yaml:"beer"`
}

type zessConfig struct {
	URL       string `yaml:"url"`
	DayAmount int    `yaml:"day_amount"`
}

type Config struct {
	Cammie  cammieConfig  `yaml:"cammie"`
	Buzzer  buzzerConfig  `yaml:"buzzer"`
	Spotify spotifyConfig `yaml:"spotify"`
	Tap     tapConfig     `yaml:"tap"`
	Zess    zessConfig    `yaml:"zess"`
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
