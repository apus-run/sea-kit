package cache

import "errors"

var (
	ErrKeyNotExist                = errors.New("key 不存在")
	ErrDeleteKeyFailed            = errors.New("删除key失败")
	ErrKeyNeverExpireNotSupported = errors.New("不支持key永不过期")
)
