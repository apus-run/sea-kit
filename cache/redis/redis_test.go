package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/apus-run/sea-kit/cache"
	"github.com/apus-run/sea-kit/cache/mocks"
)

func TestCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key        string
		value      string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "set value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")
				cmd.EXPECT().
					Set(context.Background(), "name", "小芳", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "小芳",
			expiration: time.Minute,
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Set(context.Background(), "name", "小芳", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "小芳",
			expiration: time.Minute,

			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCache_Get(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key string

		wantErr error
		wantVal string
	}{
		{
			name: "get value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetVal("小芳")
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantVal: "小芳",
		},
		{
			name: "get error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetErr(redis.Nil)
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantErr: cache.ErrKeyNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCache(tc.mock(ctrl))
			val, err := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
