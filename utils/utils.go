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

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func GetDayDifference(day1 int, date2 time.Time) int {
	daysInYear := 365
	if isLeapYear(date2.Year()) {
		daysInYear = 366
	}

	return (date2.YearDay() - day1 + daysInYear) % daysInYear
}

func ShiftSliceBackward[T any](slice []T) []T {
	newSlice := make([]T, len(slice)-1, cap(slice))
	copy(newSlice, slice[1:])

	return newSlice
}

func AddSliceElement[T any](slice []T, element T) []T {
	if len(slice) >= cap(slice) {
		// Array is max size, shift everything
		slice = ShiftSliceBackward(slice)
	}

	slice = append(slice, element)

	return slice
}
