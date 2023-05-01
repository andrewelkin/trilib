package utils

import "strings"

// SplitNoEmpty splits string and removes all empties
func SplitNoEmpty(s, sep string) (res []string) {
	for _, a := range strings.Split(s, sep) {
		if len(a) > 0 {
			res = append(res, a)
		}
	}
	return res
}

func SplitNoEmptyTrimmed(s, sep string) (res []string) {
	for _, a := range strings.Split(s, sep) {
		if len(a) > 0 {
			res = append(res, strings.Trim(a, " "))
		}
	}
	return res
}

// MakeArrayUnique makes string array unique
func MakeArrayUnique(a []string) []string {
	tmpMap := make(map[string]int)
	for i, _ := range a {
		tmpMap[a[i]] = 1
	}
	var res []string
	for k := range tmpMap {
		res = append(res, k)
	}
	return res
}
