// weighted round robin负载均衡模块，主要用于为每个RPC请求返回一个Server节点以供调用

package wrr

import (
	"sync"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "weighted_round_robin"

func init() {
	balancer.Register(newBuilder())
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(
		Name,
		&WeightedPickerBuilder{},
		base.Config{HealthCheck: true},
	)
}

type WeightedPicker struct {
	mutex sync.Mutex
	conns []*weightConn
}

func (w *WeightedPicker) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// only one connection
	if len(w.conns) == 1 {
		return balancer.PickResult{SubConn: w.conns[0].SubConn}, nil
	}

	// 这里实时计算 totalWeight 是为了方便你作业动态调整权重
	var totalWeight int
	var selected *weightConn

	w.mutex.Lock()
	for _, node := range w.conns {
		totalWeight += node.weight
		node.currentWeight += node.weight
		if selected == nil || node.currentWeight > selected.currentWeight {
			selected = node
		}
	}

	selected.currentWeight -= totalWeight
	w.mutex.Unlock()

	return balancer.PickResult{
		SubConn: selected.SubConn,
		Done: func(info balancer.DoneInfo) {
			// 在这里执行 failover 有关的事情
			// 例如说把 selected 的 currentWeight 进一步调低到一个非常低的值
			// 也可以直接把 selected 从 w.conns 删除
		},
	}, nil
}

type WeightedPickerBuilder struct{}

func (b *WeightedPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, scInfo := range info.ReadySCs {
		// 如果不存在，那么权重就是 0
		weight := 0
		if scInfo.Address.Attributes != nil {
			// weight 经过注册中心的转发之后，变成了 float64，要小心这个问题
			if w, ok := scInfo.Address.Attributes.Value("weight").(int); ok {
				weight = w
			}
		}

		conns = append(conns, &weightConn{
			weight:        weight,
			currentWeight: weight,
			addr:          scInfo.Address.Addr,

			SubConn:     sc,
			SubConnInfo: scInfo,
		})
	}

	return &WeightedPicker{
		conns: conns,
	}
}

type weightConn struct {
	// 初始权重值
	weight int
	// 当前权重值
	currentWeight int

	addr string
	balancer.SubConn
	base.SubConnInfo
}

func (w *weightConn) Address() string {
	return w.SubConnInfo.Address.Addr
}

func (w *weightConn) Attributes() *attributes.Attributes {
	return w.SubConnInfo.Address.Attributes
}
