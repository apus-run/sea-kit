package session

import (
	"context"

	"github.com/apus-run/sea-kit/collection"
)

var _ Session = &MemorySession{}

// MemorySession 一般用于测试
type MemorySession struct {
	data   map[string]any
	claims Claims
}

func (m *MemorySession) Destroy(ctx context.Context) error {
	return nil
}

func (m *MemorySession) Del(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func NewMemorySession(cl Claims) *MemorySession {
	return &MemorySession{
		data:   map[string]any{},
		claims: cl,
	}
}

func (m *MemorySession) Set(ctx context.Context, key string, val any) error {
	m.data[key] = val
	return nil
}

func (m *MemorySession) Get(ctx context.Context, key string) collection.AnyValue {
	val, ok := m.data[key]
	if !ok {
		return collection.AnyValue{Error: ErrSessionKeyNotFound}
	}
	return collection.AnyValue{Value: val}
}

func (m *MemorySession) Claims() Claims {
	return m.claims
}
