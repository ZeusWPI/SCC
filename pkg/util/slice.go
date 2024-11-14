// Package util provides utility functions
package util

import "strings"

// SliceMap maps a slice of type T to a slice of type U
func SliceMap[T any, U any](input []T, mapFunc func(T) U) []U {
	v := make([]U, len(input))
	for i, item := range input {
		v[i] = mapFunc(item)
	}
	return v
}

// SliceStringJoin joins a slice of type T to a string with a separator
func SliceStringJoin[T any](input []T, sep string, mapFunc func(T) string) string {
	v := make([]string, len(input))
	for i, item := range input {
		v[i] = mapFunc(item)
	}
	return strings.Join(v, sep)
}
