package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestFromRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromRequest(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
