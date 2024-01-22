//  P2C 算法 + EWMA 来实现自适应的 LB 机制

package p2c

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	// Name is the name of p2c balancer.
	Name = "p2c_ewma"

	decayTime       = int64(time.Second * 10) // default value from finagle
	forcePick       = int64(time.Second)
	initSuccess     = 1000
	throttleSuccess = initSuccess / 2
	penalty         = int64(math.MaxInt32)
	pickTimes       = 3
	logInterval     = time.Minute
)

var emptyPickResult balancer.PickResult

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

// Build  每次有后端节点新增/下线都会触发初始化
func (b *p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	readySCs := info.ReadySCs
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var conns []*subConn
	for conn, connInfo := range readySCs {
		conns = append(conns, &subConn{
			addr:    connInfo.Address,
			conn:    conn,
			success: initSuccess,
		})
	}

	// 每次有后端节点上下线时会调用，重新初始化 `[]*subConn`
	return &p2cPicker{
		conns: conns,
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		stamp: NewAtomicDuration(),
	}
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(p2cPickerBuilder), base.Config{HealthCheck: true})
}

type p2cPicker struct {
	conns []*subConn // 保存所有服务的节点信息
	r     *rand.Rand
	stamp *AtomicDuration
	mu    sync.Mutex
}

// Pick 会在每次请求时调用，用于选择一个节点进行请求。
// 在 Pick 方法里实现了 P2C [算法](https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/balancer/p2c/p2c.go#L75)，挑选合适的节点，并通过节点的 EWMA 值计算负载情况，返回负载低的节点供 gRPC 使用。go-zero 是使用的下面的 Pick 逻辑：
// 1.多选二，基于 P2C 算法
// 2.二再选一，基于 EWMA 负载低的原则
func (p *p2cPicker) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var chosen *subConn
	switch len(p.conns) {
	case 0:
		// 无节点
		return emptyPickResult, balancer.ErrNoSubConnAvailable
	case 1:
		// only one，直接返回
		chosen = p.choose(p.conns[0], nil)
	case 2:
		// only two，直接进入 choose，通过 EWMA 值 计算负载，并返回负载低的节点返回供 gRPC 使用
		chosen = p.choose(p.conns[0], p.conns[1])
	default:
		// 超过 3 个节点，P2C 选择两个
		var node1, node2 *subConn
		for i := 0; i < pickTimes; i++ {
			a := p.r.Intn(len(p.conns))
			b := p.r.Intn(len(p.conns) - 1)
			if b >= a {
				b++
			}
			node1 = p.conns[a]
			node2 = p.conns[b]
			// 如果这次选择的节点达到了健康要求, 就中断选择
			if node1.healthy() && node2.healthy() {
				break
			}
		}
		// 比较两个节点的负载情况，选择负载低的节点返回供 gRPC 使用
		chosen = p.choose(node1, node2)
	}

	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)

	// 根据负载均衡算法选择1个节点进行请求
	return balancer.PickResult{
		SubConn: chosen.conn,
		Done:    p.buildDoneFunc(chosen),
	}, nil
}

// buildDoneFunc
// 在请求结束的时候 gRPC 会调用 `PickResult.Done` 方法，在此实现了本次请求耗时等信息的存储，并计算出了 EWMA 值保存了起来，供下次请求时计算负载等情况的使用：
// -	被选中节点正在处理请求的总数减 `1`
// -	保存处理请求结束的时间点，用于计算距离上次节点处理请求的差值，并算出 EWMA 中的 `β` 值
// -	通过 EWMA 算法计算本次请求耗时（延迟），并保存到节点的 `lag` 属性里
// -	通过 EWMA 算法，更新计算节点的健康状态保存到节点的 `success` 属性中
func (p *p2cPicker) buildDoneFunc(c *subConn) func(info balancer.DoneInfo) {
	start := int64(Now())
	return func(info balancer.DoneInfo) {
		// 请求结束，把节点正在处理请求的总数（并发数）减 1
		atomic.AddInt64(&c.inflight, -1)
		now := Now()
		// 获取上一次请求的时间, 保存处理请求结束的时间点，返回之前的值 c.last
		last := atomic.SwapInt64(&c.last, int64(now))
		td := int64(now) - last
		if td < 0 {
			td = 0
		}

		// 计算 belta 值, 注意 td 为负, 用牛顿冷却定律中的衰减函数模型计算 EWMA 算法中的β值
		w := math.Exp(float64(-td) / float64(decayTime))
		lag := int64(now) - start
		if lag < 0 {
			lag = 0
		}
		// 保存本次请求的耗时，并返回上次的耗时，用于 ewma 计算
		olag := atomic.LoadUint64(&c.lag)
		if olag == 0 {
			w = 0
		}
		// 计算并存储 EWMA 值
		atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
		success := initSuccess
		if info.Err != nil && !Acceptable(info.Err) {
			// 非逻辑类错误，可能是超时等，需要降低节点权重
			success = 0
		}
		osucc := atomic.LoadUint64(&c.success)
		atomic.StoreUint64(&c.success, uint64(float64(osucc)*w+float64(success)*(1-w)))

		// 按需打印节点日志
		stamp := p.stamp.Load()
		if now-stamp >= logInterval {
			if p.stamp.CompareAndSwap(stamp, now) {
				p.logStats()
			}
		}
	}
}

func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	start := int64(Now())
	if c2 == nil {
		atomic.StoreInt64(&c1.pick, start)
		return c1
	}
	// 如果 c1 指向 conn 的负载比 c2 指向 conn 的负载高，那么让 c1 指向负载低的，c2 指向高的
	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	// 这段代码的用意：如果选中的节点 `c2`（相对的高负载），在 `forcePick` 期间内没有被选中一次，那么强制选择一次。这里是利用强制的机会，来触发成功率、延迟的衰减，不然可能会导致此节点永远不会被选中，比较巧妙的设计。这里还使用了 `CompareAndSwapInt64` 方法，即原子锁保证并发安全，仅放行 `1` 次。
	pick := atomic.LoadInt64(&c2.pick)
	if start-pick > forcePick && atomic.CompareAndSwapInt64(&c2.pick, pick, start) {
		return c2
	}

	//pick c1，更新 pick 时间
	atomic.StoreInt64(&c1.pick, start)
	return c1
}

func (p *p2cPicker) logStats() {
	var stats []string

	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.conns {
		stats = append(stats, fmt.Sprintf("conn: %s, load: %d, reqs: %d",
			conn.addr.Addr, conn.load(), atomic.SwapInt64(&conn.requests, 0)))
	}

	log.Printf("p2c - %s", strings.Join(stats, "; "))
}

// 加入 EMWA 值、权重、并发量等一些因子，来计算节点的负载
type subConn struct {
	lag      uint64 // 保存EWMA值
	inflight int64  // 保存当前节点正在并发的请求总数
	success  uint64 // 标识一段时间内此连接的健康状态
	requests int64  // 请求总数
	last     int64  // 上一次请求耗时, 用于计算 EWMA 值
	pick     int64  // 该连接上一次被选中的时间戳
	addr     resolver.Address
	conn     balancer.SubConn
}

// `throttleSuccess` 为常量，`c.success` 初始值为 `throttleSuccess` 的 `2` 倍
func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

// `load` 方法预估（计算）了 **节点的负载情况**，是通过下面公式实现：`load = ewma * inflight`，这里有点意思：
// >	ewma 相当于平均请求耗时
// >	inflight 是当前节点正在处理请求的数量
// 上面 `2` 个因子相乘大致计算出了当前节点的网络负载。
func (c *subConn) load() int64 {
	// 通过 EWMA 计算节点的负载情况； 加 1 是为了避免为 0 的情况
	// plus one to avoid multiply zero
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		// penalty：初始化没有数据时的惩罚值，默认为 1e9 * 250
		return penalty
	}

	return load
}
