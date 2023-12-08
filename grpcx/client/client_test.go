package client

import (
	"context"
	"google.golang.org/grpc"
	"reflect"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	o := &Options{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithDialOptions(v...)(o)
	if !reflect.DeepEqual(v, o.dialOpts) {
		t.Errorf("expect %v but got %v", v, o.dialOpts)
	}
}

func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := NewClient(
		ctx,
		WithAddr("abc"),
	)
	if err != nil {
		t.Error(err)
	}
}
