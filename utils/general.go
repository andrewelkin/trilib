package utils

import (
	"regexp"
	"sort"
	"strings"
)

const BPS_SCALAR = 10000.0

// StringInSlice checks if a string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// WildCardToRegexp
func WildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {
		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}
		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

// MatchMask
func MatchMask(pattern string, value string) bool {
	result, _ := regexp.MatchString(WildCardToRegexp(pattern), value)
	return result
}

// StringInSliceLastToFirst checks if a string is in a slice going from last element to first
func StringInSliceLastToFirst(a string, list []string) bool {
	for i := len(list) - 1; i >= 0; i-- {
		if list[i] == a {
			return true
		}

	}
	return false
}

// BuyQToSide converts buyQ boolean into side string
func BuyQToSide(buyQ bool) string {
	var side string
	if buyQ {
		side = "buy"
	} else {
		side = "sell"
	}
	return side
}

// MaxIntegerInSlice returns the largest interger in a slice of integers
func MaxIntegerInSlice(intSlice []int) int {
	max := intSlice[0]
	for _, value := range intSlice {
		if value > max {
			max = value
		}
	}
	return max
}

// Provides map value if one exists for given key otherwise sets a default value
func SetDefault(h map[string]int, k string, v int) (set bool, r int) {
	if r, set = h[k]; !set {
		h[k] = v
		r = v
		set = true
	}
	return
}

func RemoveStringInSlice(slice []string, stringToRemove string) (parsedSlice []string) {
	if StringInSlice(stringToRemove, slice) {
		for index, s := range slice {
			if s == stringToRemove {
				parsedSlice = append(slice[:index], slice[index+1:]...)
				break
			}
		}
	} else {
		parsedSlice = slice
	}
	return
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// ReverseSlice reverses slice in place
func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

// ReversedSlice returns a copy of a slice
func ReversedSlice(s []interface{}) []interface{} {
	a := make([]interface{}, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func MakeSureNegativeIf[T float64 | int64 | float32 | int | int32](condition bool, v T) T {
	if condition {
		if v > 0 {
			v = -v
		}
	} else {
		if v < 0 {
			v = -v
		}
	}
	return v
}

func GetValueOrDefault[K comparable, V comparable](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}
