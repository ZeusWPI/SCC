package utils

import (
	"math/rand/v2"
	"time"
)

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func TimeAndDateFormat() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04 02/01")
	return formattedTime
}
