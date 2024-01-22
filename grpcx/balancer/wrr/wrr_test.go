package wrr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func TestWrrPicker_PickNil(t *testing.T) {
	builder := new(WeightedPickerBuilder)
	picker := builder.Build(base.PickerBuildInfo{})
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.NotNil(t, err)
}

func TestPickerWithEmptyConns(t *testing.T) {
	var picker = &WeightedPicker{}
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.ErrorIs(t, err, balancer.ErrNoSubConnAvailable)
}
