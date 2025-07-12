package empty

import (
	"context"
)

type NoCopy struct{}

func (*NoCopy) Lock(ctx context.Context) error {
	return nil
}

func (*NoCopy) Unlock(ctx context.Context) error {
	return nil
}

func (*NoCopy) Refresh(ctx context.Context) error {
	return nil
}
