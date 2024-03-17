package session

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apus-run/sea-kit/collection"
)

func TestMemorySession_GetSet(t *testing.T) {
	testCases := []struct {
		name string

		// 插入数据
		key string
		val string

		getKey string

		wantVal collection.AnyValue
	}{
		{
			name:    "成功获取",
			key:     "key1",
			val:     "value1",
			getKey:  "key1",
			wantVal: collection.AnyValue{Value: "value1"},
		},
		{
			name:    "没有数据",
			key:     "key1",
			val:     "value1",
			getKey:  "key2",
			wantVal: collection.AnyValue{Error: ErrSessionKeyNotFound},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := NewMemorySession(Claims{})
			ctx := context.Background()
			err := ms.Set(ctx, tc.key, tc.val)
			require.NoError(t, err)
			val := ms.Get(ctx, tc.getKey)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
