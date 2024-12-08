// Package config provides all configuration related functions
package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// FIXME: Add mutex for map writes

var mu sync.Mutex

func bindEnv(key string) {
	envName := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

	mu.Lock()
	defer mu.Unlock()

	viper.BindEnv(key, envName)
}

// Init initializes the configuration
func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	viper.AutomaticEnv()
	env := GetDefaultString("app.env", "development")

	viper.SetConfigName(fmt.Sprintf("%s.yaml", strings.ToLower(env)))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	return viper.ReadInConfig()
}

// GetString returns the value of the key in string
func GetString(key string) string {
	bindEnv(key)

	mu.Lock()
	defer mu.Unlock()

	return viper.GetString(key)
}

// GetDefaultString returns the value of the key in string or a default value
func GetDefaultString(key, defaultValue string) string {
	mu.Lock()
	defer mu.Unlock()

	viper.SetDefault(key, defaultValue)
	return GetString(key)
}

// GetStringSlice returns the value of the key in string slice
func GetStringSlice(key string) []string {
	bindEnv(key)

	mu.Lock()
	defer mu.Unlock()

	return viper.GetStringSlice(key)
}

// GetDefaultStringSlice returns the value of the key in string slice or a default value
func GetDefaultStringSlice(key string, defaultValue []string) []string {
	mu.Lock()
	defer mu.Unlock()

	viper.SetDefault(key, defaultValue)
	return GetStringSlice(key)
}

// GetInt returns the value of the key in int
func GetInt(key string) int {
	bindEnv(key)

	mu.Lock()
	defer mu.Unlock()

	return viper.GetInt(key)
}

// GetDefaultInt returns the value of the key in int or a default value
func GetDefaultInt(key string, defaultVal int) int {
	mu.Lock()
	defer mu.Unlock()

	viper.SetDefault(key, defaultVal)
	return GetInt(key)
}

// GetUint16 returns the value of the key in uint16
func GetUint16(key string) uint16 {
	bindEnv(key)

	mu.Lock()
	defer mu.Unlock()

	return viper.GetUint16(key)
}

// GetDefaultUint16 returns the value of the key in uint16 or a default value
func GetDefaultUint16(key string, defaultVal uint16) uint16 {
	mu.Lock()
	defer mu.Unlock()

	viper.SetDefault(key, defaultVal)
	return GetUint16(key)
}

// GetBool returns the value of the key in bool
func GetBool(key string) bool {
	bindEnv(key)

	mu.Lock()
	defer mu.Unlock()

	return viper.GetBool(key)
}

// GetDefaultBool returns the value of the key in bool or a default value
func GetDefaultBool(key string, defaultVal bool) bool {
	mu.Lock()
	defer mu.Unlock()

	viper.SetDefault(key, defaultVal)
	return GetBool(key)
}
