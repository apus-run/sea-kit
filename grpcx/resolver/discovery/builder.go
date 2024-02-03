package discovery

import (
	log "github.com/apus-run/sea-kit/zlog"
	"google.golang.org/grpc/resolver"
)

// Option is builder option.
type Option func(o *discoveryBuilder)

// WithLogger with logger option.
func WithLogger(log log.Logger) Option {
	return func(b *discoveryBuilder) {
		b.logger = log
	}
}

func WithDiscoveryMinPeers(n int) Option {
	return func(b *discoveryBuilder) {
		b.discoveryMinPeers = n
	}
}

func WithDiscoCh(c chan []string) Option {
	return func(b *discoveryBuilder) {
		b.discoCh = c
	}
}

// PrintDebugLog print grpc resolver watch service log
func PrintDebugLog(p bool) Option {
	return func(b *discoveryBuilder) {
		b.debugLog = p
	}
}

type discoveryBuilder struct {
	*discoveryResolver
}

// NewBuilder creates a builder
func NewBuilder(notifier Notifier, discoverer Discoverer, opts ...Option) resolver.Builder {
	b := &discoveryBuilder{
		New(notifier, discoverer, 3, make(chan []string, 100), log.L(), true),
	}
	for _, o := range opts {
		o(b)
	}

	return b
}

// Build creates a new resolver
func (b *discoveryBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return b.discoveryResolver.Build(target, cc, opts)
}

// Scheme return scheme of discovery
func (b *discoveryBuilder) Scheme() string {
	return b.scheme
}
