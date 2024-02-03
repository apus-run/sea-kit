package discovery

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	log "github.com/apus-run/sea-kit/zlog"
	"google.golang.org/grpc/resolver"
)

// GRPCServiceConfig provides grpc service config
const GRPCServiceConfig = `{"loadBalancingPolicy":"round_robin"}`

// Resolver uses notifier to fetch list of available hosts
type discoveryResolver struct {
	scheme            string
	cc                resolver.ClientConn
	notifier          Notifier
	discoverer        Discoverer
	logger            log.Logger
	debugLog          bool
	discoCh           chan []string // used to receive notifications
	discoveryMinPeers int
	salt              []byte

	// used to block Close() until the watcher goroutine exits its loop
	closing sync.WaitGroup
}
type hostScore struct {
	address string
	score   uint32
}

type hostScores []hostScore

func (s hostScores) Len() int           { return len(s) }
func (s hostScores) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s hostScores) Less(i, j int) bool { return s[i].score < s[j].score }

// New initialize a new grpc resolver with notifier
func New(
	notifier Notifier,
	discoverer Discoverer,
	discoveryMinPeers int,
	discoCh chan []string,
	logger log.Logger,
	debugLog bool,
) *discoveryResolver {
	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))
	r := &discoveryResolver{
		notifier:          notifier,
		discoverer:        discoverer,
		discoCh:           discoCh,
		logger:            logger,
		debugLog:          debugLog,
		discoveryMinPeers: discoveryMinPeers,
		salt:              []byte(strconv.FormatInt(random.Int63(), 10)), // random salt for rendezvousHash
		scheme:            strconv.FormatInt(seed, 36),                   // make random scheme which will be used when registering
	}

	// Register the resolver with grpc so it's available for grpc.Dial
	resolver.Register(r)

	// Register the discoCh channel with notifier so it continues to fetch a list of host/port
	notifier.Register(r.discoCh)
	return r
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (r *discoveryResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	// Update conn states if proactively updates already work
	instances, err := r.discoverer.Instances()
	if err != nil {
		return nil, err
	}
	r.updateAddresses(instances)
	r.closing.Add(1)
	go r.watcher()
	return r, nil
}

// Scheme returns resolver's scheme.
func (r *discoveryResolver) Scheme() string {
	return r.scheme
}

// ResolveNow is a noop for Resolver since resolver is already firing r.cc.UpdatesState every time
// it receives updates of new instance from discoCh
func (r *discoveryResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (r *discoveryResolver) watcher() {
	defer r.closing.Done()
	for latestHostPorts := range r.discoCh {
		if r.debugLog {
			r.logger.Info("Received updates from notifier", log.Strings("hostPorts", latestHostPorts))
		}
		r.updateAddresses(latestHostPorts)
	}
}

// Close closes both discoCh
func (r *discoveryResolver) Close() {
	r.notifier.Unregister(r.discoCh)
	close(r.discoCh)
	r.closing.Wait()
}

// rendezvousHash is the core of the algorithm. It takes input addresses,
// assigns each of them a hash, sorts them by those hash values, and
// returns top N of entries from the sorted list, up to minPeers parameter.
func (r *discoveryResolver) rendezvousHash(addresses []string) []string {
	hasher := fnv.New32()
	hosts := make(hostScores, len(addresses))
	for i, address := range addresses {
		hosts[i] = hostScore{
			address: address,
			score:   hashAddr(hasher, []byte(address), r.salt),
		}
	}
	sort.Sort(hosts)
	n := min(r.discoveryMinPeers, len(hosts))
	topN := make([]string, n)
	for i := 0; i < n; i++ {
		topN[i] = hosts[i].address
	}
	return topN
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func hashAddr(hasher hash.Hash32, node, saltKey []byte) uint32 {
	hasher.Reset()
	hasher.Write(saltKey)
	hasher.Write(node)
	return hasher.Sum32()
}

func (r *discoveryResolver) updateAddresses(hostPorts []string) {
	topN := r.rendezvousHash(hostPorts)
	addresses := generateAddresses(topN)
	r.cc.UpdateState(resolver.State{Addresses: addresses})
}

func generateAddresses(instances []string) []resolver.Address {
	addrs := make([]resolver.Address, len(instances))
	for i, instance := range instances {
		addrs[i] = resolver.Address{Addr: instance}
	}
	return addrs
}
