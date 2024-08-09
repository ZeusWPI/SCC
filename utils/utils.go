package utils

import (
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func GetEnvAsBool(name string, defaultValue bool) bool {
	valueStr := os.Getenv(name)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func TimeAndDateFormat() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04 02/01")
	return formattedTime
}
