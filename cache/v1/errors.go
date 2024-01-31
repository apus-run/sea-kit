package cache

import "errors"

var (
	ErrCacheExpired = errors.New("tiny_cache expired")
	ErrKeyNotFound  = errors.New("key not found")
	ErrTypeNotOk    = errors.New("val type not ok")
)
