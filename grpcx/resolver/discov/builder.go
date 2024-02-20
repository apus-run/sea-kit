package discov

import (
	"context"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"

	"github.com/apus-run/sea-kit/grpcx/registry"
	log "github.com/apus-run/sea-kit/zlog"
)

const name = "discovery"

// Option is builder option.
type Option func(o *discovBuilder)

// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(b *discovBuilder) {
		b.timeout = timeout
	}
}

// WithInsecure with isSecure option.
func WithInsecure(insecure bool) Option {
	return func(b *discovBuilder) {
		b.insecure = insecure
	}
}

// WithLogger with logger option.
func WithLogger(log log.Logger) Option {
	return func(b *discovBuilder) {
		b.log = log
	}
}

// PrintDebugLog print grpc resolver watch service log
func PrintDebugLog(p bool) Option {
	return func(b *discovBuilder) {
		b.debugLog = p
	}
}

type discovBuilder struct {
	discoverer registry.Discovery
	timeout    time.Duration
	insecure   bool

	log      log.Logger
	debugLog bool
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &discovBuilder{
		discoverer: d,
		timeout:    time.Second * 10,
		insecure:   false,
		debugLog:   true,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	watchRes := &struct {
		err error
		w   registry.Watcher
	}{}

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		w, err := b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		watchRes.w = w
		watchRes.err = err
		close(done)
	}()

	var err error
	select {
	case <-done:
		err = watchRes.err
	case <-time.After(b.timeout):
		err = errors.New("discov create watcher overtime")
	}
	if err != nil {
		cancel()
		return nil, err
	}

	r := &discovResolver{
		w:        watchRes.w,
		cc:       cc,
		ctx:      ctx,
		cancel:   cancel,
		insecure: b.insecure,
		debugLog: b.debugLog,
		log:      b.log,
	}
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (b *discovBuilder) Scheme() string {
	return name
}
