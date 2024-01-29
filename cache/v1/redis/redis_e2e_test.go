//go:build e2e

package redis

import (
	"context"
	"encoding/json"

	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Bar struct {
	Name string
}

func (b *Bar) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Bar) UnmarshalBinary(bt []byte) error {
	return json.Unmarshal(bt, b)
}

func TestCache_All(t *testing.T) {
	Convey("test get client", t, func() {
		rdb := redis.NewClient(&redis.Options{
			Password: "123456",
			DB:       1,
			Addr:     "localhost:16379",
		})
		require.NoError(t, rdb.Ping(context.Background()).Err())
		mc := NewCache(rdb)

		So(mc, ShouldNotBeNil)
		ctx := context.Background()

		Convey("string get set", func() {
			err := mc.Set(ctx, "foo", "bar", 1*time.Hour)
			So(err, ShouldBeNil)
			val, err := mc.Get(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "bar")
			err = mc.SetTTL(ctx, "foo", 1*time.Minute)
			So(err, ShouldBeNil)
			du, err := mc.GetTTL(ctx, "foo")
			So(err, ShouldBeNil)
			So(du, ShouldBeLessThanOrEqualTo, 1*time.Minute)
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
			val, err = mc.Get(ctx, "foo")
			So(err, ShouldEqual, v1.ErrKeyNotFound)
		})

		Convey("obj get set", func() {
			obj := &Bar{
				Name: "bar",
			}
			err := mc.SetObj(ctx, "foo", obj, 1*time.Hour)
			So(err, ShouldBeNil)
			objNew := Bar{}
			err = mc.GetObj(ctx, "foo", &objNew)
			So(err, ShouldBeNil)
			So(objNew.Name, ShouldEqual, "bar")
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
		})

		Convey("many op", func() {
			err := mc.SetMany(ctx, map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			}, 1*time.Hour)
			So(err, ShouldBeNil)

			ret, err := mc.GetMany(ctx, []string{"foo1", "foo2"})
			So(err, ShouldBeNil)
			So(len(ret), ShouldEqual, 2)
			So(ret, ShouldContainKey, "foo2")
			So(ret["foo2"], ShouldEqual, "bar2")

			err = mc.DelMany(ctx, []string{"foo1", "foo2"})
			So(err, ShouldBeNil)
		})

		Convey("calc op", func() {
			val, err := mc.Increment(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1)
			val, err = mc.Calc(ctx, "foo", 2)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3)
			val, err = mc.Decrement(ctx, "foo")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2)
			err = mc.Del(ctx, "foo")
			So(err, ShouldBeNil)
		})
	})
}

func TestCache_e2e_Set(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Password: "123456",
		DB:       1,
		Addr:     "localhost:16379",
	})
	require.NoError(t, rdb.Ping(context.Background()).Err())

	testCases := []struct {
		name  string
		after func(ctx context.Context, t *testing.T)

		key        string
		val        string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "set e2e value",
			after: func(ctx context.Context, t *testing.T) {
				result, err := rdb.Get(ctx, "name").Result()
				require.NoError(t, err)
				assert.Equal(t, "小芳", result)

				_, err = rdb.Del(ctx, "name").Result()
				require.NoError(t, err)
			},
			key:        "name",
			val:        "小芳",
			expiration: time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
			defer cancelFunc()
			c := NewCache(rdb)

			err := c.Set(ctx, "name", "小芳", time.Minute)
			assert.NoError(t, err)
			tc.after(ctx, t)
		})
	}
}

func TestCache_e2e_Get(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	require.NoError(t, rdb.Ping(context.Background()).Err())

	testCases := []struct {
		name   string
		before func(ctx context.Context, t *testing.T)
		after  func(ctx context.Context, t *testing.T)

		key string

		wantVal string
		wantErr error
	}{
		{
			name: "get e2e value",
			before: func(ctx context.Context, t *testing.T) {
				require.NoError(t, rdb.Set(ctx, "name", "小芳", time.Minute).Err())
			},
			after: func(ctx context.Context, t *testing.T) {
				require.NoError(t, rdb.Del(ctx, "name").Err())
			},
			key: "name",

			wantVal: "小芳",
		},
		{
			name:    "get e2e error",
			key:     "name",
			before:  func(ctx context.Context, t *testing.T) {},
			after:   func(ctx context.Context, t *testing.T) {},
			wantErr: v1.ErrKeyNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
			defer cancelFunc()
			c := NewCache(rdb)

			tc.before(ctx, t)
			val, err := c.Get(ctx, tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
			tc.after(ctx, t)
		})
	}
}
