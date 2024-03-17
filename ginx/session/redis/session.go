package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/apus-run/sea-kit/collection"
	"github.com/apus-run/sea-kit/ginx/session"
)

var _ session.Session = &Session{}

// Session 生命周期应该和 http 请求保持一致
type Session struct {
	client redis.Cmdable
	// key 是 ssid 拼接而成。注意，它不是 access token，也不是 refresh token
	key        string
	claims     session.Claims
	expiration time.Duration
}

func (sess *Session) Destroy(ctx context.Context) error {
	return sess.client.Del(ctx, sess.key).Err()
}

func (sess *Session) Del(ctx context.Context, key string) error {
	return sess.client.Del(ctx, sess.key, key).Err()
}

func (sess *Session) Set(ctx context.Context, key string, val any) error {
	return sess.client.HSet(ctx, sess.key, key, val).Err()
}

func (sess *Session) init(ctx context.Context, kvs map[string]any) error {
	pip := sess.client.Pipeline()
	for k, v := range kvs {
		pip.HMSet(ctx, sess.key, k, v)
	}
	pip.Expire(ctx, sess.key, sess.expiration)
	_, err := pip.Exec(ctx)
	return err
}

func (sess *Session) Get(ctx context.Context, key string) collection.AnyValue {
	res, err := sess.client.HGet(ctx, sess.key, key).Result()
	if err != nil {
		return collection.AnyValue{Error: err}
	}
	return collection.AnyValue{
		Value: res,
	}
}

func (sess *Session) Claims() session.Claims {
	return sess.claims
}

func newRedisSession(
	ssid string,
	expiration time.Duration,
	client redis.Cmdable, cl session.Claims) *Session {
	return &Session{
		client:     client,
		key:        "session:" + ssid,
		expiration: expiration,
		claims:     cl,
	}
}
