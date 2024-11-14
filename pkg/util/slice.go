// Package util provides utility functions
package util

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
	str := ""
	for _, item := range input {
		str += mapFunc(item) + sep
	}
	return str[:len(str)-len(sep)]
}

// SliceFilter filters a slice of type T based on a filter function
func SliceFilter[T any](input []T, filterFunc func(T) bool) []T {
	v := make([]T, 0)
	for _, item := range input {
		if filterFunc(item) {
			v = append(v, item)
		}
	}
	return v
}
