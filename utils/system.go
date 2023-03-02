package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(0)
	if !ok {
		return "unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unnamed"
	}
	return me.Name()
}

func GetFunctionName(f interface{}, short bool) string {
	n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if short {
		if nn := strings.LastIndex(n, "/"); nn > 0 {
			n = n[nn:]
		}
	}
	return n
}
