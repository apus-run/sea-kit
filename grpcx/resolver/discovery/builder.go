package discovery

import (
	"math/rand"
	"strconv"
	"time"

	log "github.com/apus-run/sea-kit/zlog"
	"google.golang.org/grpc/resolver"
)

// Option is builder option.
type Option func(o *discoveryBuilder)

// WithLogger with logger option.
func WithLogger(log log.Logger) Option {
	return func(b *discoveryBuilder) {
		b.log = log
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
	scheme            string
	notifier          Notifier
	discoverer        Discoverer
	discoCh           chan []string // used to receive notifications
	discoveryMinPeers int
	salt              []byte

	log      log.Logger
	debugLog bool
}

// NewBuilder creates a builder
func NewBuilder(notifier Notifier, discoverer Discoverer, opts ...Option) resolver.Builder {
	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	b := &discoveryBuilder{
		scheme:            strconv.FormatInt(seed, 36),
		notifier:          notifier,
		discoverer:        discoverer,
		debugLog:          true,
		discoveryMinPeers: 3,
		discoCh:           make(chan []string, 100),
		// random salt for rendezvousHash
		salt: []byte(strconv.FormatInt(random.Int63(), 10)),
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (b *discoveryBuilder) Build(_ resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	r := &discoveryResolver{
		cc:                cc,
		notifier:          b.notifier,
		discoverer:        b.discoverer,
		discoCh:           b.discoCh,
		log:               b.log,
		discoveryMinPeers: b.discoveryMinPeers,
		salt:              b.salt,
	}

	// Register the resolver with grpc so it's available for grpc.Dial
	resolver.Register(b)

	// Register the discoCh channel with notifier so it continues to fetch a list of host/port
	b.notifier.Register(b.discoCh)
	// Update conn states if proactively updates already work
	instances, err := b.discoverer.Instances()
	if err != nil {
		return nil, err
	}
	r.updateAddresses(instances)
	r.closing.Add(1)
	go r.watcher()

	return r, nil
}

// Scheme return scheme of discovery
func (b *discoveryBuilder) Scheme() string {
	return b.scheme
}
