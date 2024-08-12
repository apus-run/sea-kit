package hooks

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	// 创建一个新的 Hooker 实例
	hooker := New()

	hooker.Add(HookFunc{
		Name: "One",
		Fn: func(ctx context.Context) error {
			t.Log("One")
			return nil
		},
	})

	hooker.Add(HookFunc{
		Name: "Two",
		Fn: func(ctx context.Context) error {
			t.Log("Two")
			return nil
		},
	})

	hooker.Add(HookFunc{
		Name: "Three",
		Fn: func(ctx context.Context) error {
			t.Log("Three")
			return nil
		},
	})

	// 执行 Hooker 实例
	defer hooker.Do(context.Background())
}
