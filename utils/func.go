package utils

import (
	"reflect"
	"runtime"
)

// GetFuncName return the function name
func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
