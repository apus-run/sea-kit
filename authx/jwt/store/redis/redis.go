package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store redis storage.
type Store struct {
	client redis.Cmdable

	prefix string
}

// NewStore create an *Store instance to handle token storage, deletion, and checking.
func NewStore(client redis.Cmdable, prefix string) *Store {
	return &Store{client: client, prefix: prefix}
}

// Set call the Redis client to set a key-value pair with an
// expiration time, where the key name format is <prefix><accessToken>.
func (s *Store) Set(ctx context.Context, accessToken string, val any, expiration time.Duration) error {
	cmd := s.client.Set(ctx, s.key(accessToken), val, expiration)
	return cmd.Err()
}

// Delete delete the specified JWT Token in Redis.
func (s *Store) Delete(ctx context.Context, accessToken string) (bool, error) {
	cmd := s.client.Del(ctx, s.key(accessToken))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// Check check if the specified JWT Token exists in Redis.
func (s *Store) Check(ctx context.Context, accessToken string) (bool, error) {
	s.client.Get(ctx, s.key(accessToken))

	cmd := s.client.Exists(ctx, s.key(accessToken))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// wrapperKey is used to build the key name in Redis.
func (s *Store) key(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}
