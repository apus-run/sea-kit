package utils

import (
	"reflect"
	"testing"
)

func TestGetFuncName(t *testing.T) {
	got := GetFuncName(GetFuncName)
	want := "github.com/apus-run/sea-kit/utils.GetFuncName"
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
