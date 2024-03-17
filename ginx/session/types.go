package session

import (
	"context"
	"errors"

	"github.com/apus-run/sea-kit/collection"
	"github.com/apus-run/sea-kit/ginx"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrSessionKeyNotFound = errors.New("session key not found")

// Session 混合了 JWT 的设计。
type Session interface {
	// Set 将数据写入到 Session 里面
	Set(ctx context.Context, key string, val any) error
	// Get 从 Session 中获取数据，注意，这个方法不会从 JWT 里面获取数据
	Get(ctx context.Context, key string) collection.AnyValue
	// Del 删除对应的数据
	Del(ctx context.Context, key string) error
	// Destroy 销毁整个 Session
	Destroy(ctx context.Context) error
	// Claims 编码进去了 JWT 里面的数据
	Claims() Claims
}

// Provider 定义了 Session 的整个管理机制。
// 所有的 Session 都必须支持 jwt
type Provider interface {
	// NewSession 将会初始化 Session
	// 其中 jwtData 将编码进去 jwt 中
	// sessData 将被放进去 Session 中
	NewSession(ctx *ginx.Context, uid int64, jwtData map[string]string,
		sessData map[string]any) (Session, error)
	// Get 尝试拿到 Session，如果没有，返回 error
	// Get 必须校验 Session 的合法性。
	// 也就是，用户可以预期拿到的 Session 永远是没有过期，直接可用的
	Get(ctx *ginx.Context) (Session, error)

	// RenewAccessToken 刷新并且返回一个新的 access token
	// 这个过程会校验长 token 的合法性
	RenewAccessToken(ctx *ginx.Context) error
}

type Claims struct {
	Uid  int64
	SSID string
	Data map[string]string
}

func (c Claims) Get(key string) collection.AnyValue {
	val, ok := c.Data[key]
	if !ok {
		return collection.AnyValue{Error: ErrSessionKeyNotFound}
	}
	return collection.AnyValue{Value: val}
}
