package redis

import (
	"context"
	"testing"
	"time"

	"github.com/apus-run/sea-kit/cache/v2/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/apus-run/sea-kit/cache/v2"
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
					Set(context.Background(), "name", "foo", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "foo",
			expiration: time.Minute,
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Set(context.Background(), "name", "foo", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "foo",
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
				status.SetVal("foo")
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantVal: "foo",
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

			wantErr: cache.ErrKeyNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCache(tc.mock(ctrl))
			val := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, val.Error)
			if val.Error != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val.Value.(string))
		})
	}
}

func TestCache_SetNX(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(*gomock.Controller) redis.Cmdable
		key        string
		val        string
		expiration time.Duration
		result     bool
	}{
		{
			name: "setnx value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(true)
				cmd.EXPECT().
					SetNX(context.Background(), "setnx_key", "hello tiny_cache", time.Second*10).
					Return(boolCmd)
				return cmd
			},
			key:        "setnx_key",
			val:        "hello tiny_cache",
			expiration: time.Second * 10,
			result:     true,
		},
		{
			name: "setnx error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(false)
				cmd.EXPECT().
					SetNX(context.Background(), "setnx-key", "hello tiny_cache", time.Second*10).
					Return(boolCmd)

				return cmd
			},
			key:        "setnx-key",
			val:        "hello tiny_cache",
			expiration: time.Second * 10,
			result:     false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			val, err := c.SetNX(context.Background(), tc.key, tc.val, tc.expiration)
			require.NoError(t, err)
			assert.Equal(t, tc.result, val)
		})
	}
}

func TestCache_GetSet(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     string
		wantErr error
	}{
		{
			name: "getset value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				str := redis.NewStringCmd(context.Background())
				str.SetVal("hello tiny_cache")
				cmd.EXPECT().
					GetSet(context.Background(), "test_get_set", "hello go").
					Return(str)
				return cmd
			},
			key: "test_get_set",
			val: "hello go",
		},
		{
			name: "getset error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				str := redis.NewStringCmd(context.Background())
				str.SetErr(redis.Nil)
				cmd.EXPECT().
					GetSet(context.Background(), "test_get_set_err", "hello tiny_cache").
					Return(str)
				return cmd
			},
			key:     "test_get_set_err",
			val:     "hello tiny_cache",
			wantErr: cache.ErrKeyNotExist,
		},
	}

	for _, tc := range testCase {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := NewCache(tc.mock(ctrl))
		val := c.GetSet(context.Background(), tc.key, tc.val)
		assert.Equal(t, tc.wantErr, val.Error)
	}
}

func TestCache_Delete(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key []string

		wantN   int64
		wantErr error
	}{
		{
			name: "delete single existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(1))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   []string{"name"},
			wantN: 1,
		},
		{
			name: "delete single does not existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(0))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key: []string{"name"},
		},
		{
			name: "delete multiple existed keys",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(2))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   []string{"name", "age"},
			wantN: 2,
		},
		{
			name: "delete multiple do not existed keys",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(0)
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key: []string{"name", "age"},
		},
		{
			name: "delete multiple keys, some do not existed keys",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(1)
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   []string{"name", "age", "addr"},
			wantN: 1,
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(0)
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:     []string{"name"},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCache(tc.mock(ctrl))
			n, err := c.Delete(context.Background(), tc.key...)
			assert.Equal(t, tc.wantN, n)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCache_LPush(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     []any
		wantVal int64
		wantErr error
	}{
		{
			name: "lpush value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(2)
				cmd.EXPECT().
					LPush(context.Background(), "test_list_push", "1", "2").
					Return(result)
				return cmd
			},
			key:     "test_list_push",
			val:     []any{"1", "2"},
			wantVal: 2,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			length, err := c.LPush(context.Background(), tc.key, tc.val...)
			assert.Equal(t, tc.wantVal, length)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCache_LPop(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		wantVal string
		wantErr error
	}{
		{
			name: "lpop value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				str := redis.NewStringCmd(context.Background())
				str.SetVal("test")
				cmd.EXPECT().
					LPop(context.Background(), "test_cache_lpop").
					Return(str)
				return cmd
			},
			key:     "test_cache_lpop",
			wantVal: "test",
		},
		{
			name: "lpop error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				str := redis.NewStringCmd(context.Background())
				str.SetErr(redis.Nil)
				cmd.EXPECT().
					LPop(context.Background(), "test_cache_lpop").
					Return(str)
				return cmd
			},
			key:     "test_cache_lpop",
			wantVal: "",
			wantErr: cache.ErrKeyNotExist,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			val := c.LPop(context.Background(), tc.key)
			assert.Equal(t, tc.wantVal, val.Value)
			assert.Equal(t, tc.wantErr, val.Error)
		})
	}
}

func TestCache_SAdd(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     []any
		wantVal int64
		wantErr error
	}{
		{
			name: "sadd value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(2)
				cmd.EXPECT().
					SAdd(context.Background(), "test_sadd", "hello ", "hello go").
					Return(result)
				return cmd
			},
			key:     "test_sadd",
			val:     []any{"hello", "hello go"},
			wantVal: 2,
		},
		{
			name: "sadd ignore",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(1)
				cmd.EXPECT().
					SAdd(context.Background(), "test_sadd", "hello", "hello").
					Return(result)
				return cmd
			},
			key:     "test_sadd",
			val:     []any{"hello", "hello"},
			wantVal: 1,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			length, err := c.SAdd(context.Background(), tc.key, tc.val...)
			assert.Equal(t, length, tc.wantVal)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestCache_SRem(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     []any
		wantVal int64
		wantErr error
	}{
		{
			name: "srem value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(2)
				cmd.EXPECT().
					SRem(context.Background(), "test_srem", "hello", "hello go").
					Return(result)
				return cmd
			},
			key:     "test_srem",
			val:     []any{"hello", "hello go"},
			wantVal: 2,
		},
		{
			name: "srem ignore",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(0)
				cmd.EXPECT().
					SRem(context.Background(), "test_srem", "hello").
					Return(result)
				return cmd
			},
			key:     "test_srem",
			val:     []any{"hello"},
			wantVal: 0,
		},
		{
			name: "srem error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(0)
				result.SetErr(nil)
				cmd.EXPECT().
					SRem(context.Background(), "test_srem", "hello").
					Return(result)
				return cmd
			},
			key:     "test_srem",
			val:     []any{"hello"},
			wantVal: 0,
			wantErr: nil,
		},
		{
			name: "srem section ignore",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(1)
				cmd.EXPECT().
					SRem(context.Background(), "test_srem", "hello", "go").
					Return(result)
				return cmd
			},
			key:     "test_srem",
			val:     []any{"hello", "go"},
			wantVal: 1,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			result, err := c.SRem(context.Background(), tc.key, tc.val...)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, result)
		})
	}
}

func TestCache_IncrBy(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     int64
		wantVal int64
		wantErr error
	}{
		{
			name: "tiny_cache incr",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(int64(1))
				cmd.EXPECT().
					IncrBy(context.Background(), "test_incr", int64(1)).
					Return(result)
				return cmd
			},
			key:     "test_incr",
			val:     1,
			wantVal: 1,
		},
		{
			name: "tiny_cache incr not zero",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(int64(21))
				cmd.EXPECT().
					IncrBy(context.Background(), "test_incr", int64(20)).
					Return(result)
				return cmd
			},
			key:     "test_incr",
			val:     20,
			wantVal: 21,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			result, err := c.IncrBy(context.Background(), tc.key, tc.val)
			assert.Equal(t, result, tc.wantVal)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestCache_DecrBy(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     int64
		wantVal int64
		wantErr error
	}{
		{
			name: "tiny_cache decr",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(int64(0))
				cmd.EXPECT().
					DecrBy(context.Background(), "test_cache_decr", int64(1)).
					Return(result)
				return cmd
			},
			key: "test_cache_decr",
			val: 1,
		},
		{
			name: "tiny_cache decr not zero",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(int64(10))
				cmd.EXPECT().
					DecrBy(context.Background(), "test_cache_decr", int64(20)).
					Return(result)
				return cmd
			},
			key:     "test_cache_decr",
			val:     20,
			wantVal: 10,
		},
		{
			name: "tiny_cache decr negative number",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewIntCmd(context.Background())
				result.SetVal(int64(-1))
				cmd.EXPECT().
					DecrBy(context.Background(), "test_cache_decr", int64(1)).
					Return(result)
				return cmd
			},
			key:     "test_cache_decr",
			val:     1,
			wantVal: -1,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			result, err := c.DecrBy(context.Background(), tc.key, tc.val)
			assert.Equal(t, result, tc.wantVal)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestCache_IncrByFloat(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(*gomock.Controller) redis.Cmdable
		key     string
		val     float64
		wantVal float64
		wantErr error
	}{
		{
			name: "tiny_cache incrbyfloat",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewFloatCmd(context.Background())
				result.SetVal(1.2)
				cmd.EXPECT().
					IncrByFloat(context.Background(), "test_cache_incrbyfloat", 1.2).
					Return(result)

				return cmd
			},
			key:     "test_cache_incrbyfloat",
			val:     1.2,
			wantVal: 1.2,
		},
		{
			name: "tiny_cache incrbyfloat decr value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewFloatCmd(context.Background())
				result.SetVal(float64(-2.0))
				cmd.EXPECT().
					IncrByFloat(context.Background(), "test_cache_incrbyfloat", -1.0).
					Return(result)

				return cmd
			},
			key:     "test_cache_incrbyfloat",
			val:     -1,
			wantVal: -2,
		},
		{
			name: "tiny_cache incrbyfloat zero value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewFloatCmd(context.Background())
				result.SetVal(0.0)
				cmd.EXPECT().
					IncrByFloat(context.Background(), "test_cache_incrbyfloat", -12.0).
					Return(result)

				return cmd
			},
			key:     "test_cache_incrbyfloat",
			val:     -12.0,
			wantVal: 0.0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCache(tc.mock(ctrl))
			result, err := c.IncrByFloat(context.Background(), tc.key, tc.val)
			assert.Equal(t, tc.wantVal, result)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
