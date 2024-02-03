package discov

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/apus-run/sea-kit/grpcx/internal/endpoint"
	"github.com/apus-run/sea-kit/grpcx/registry"
	log "github.com/apus-run/sea-kit/zlog"
)

type discovResolver struct {
	w  registry.Watcher
	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure bool
	debugLog bool
	log      log.Logger
}

func (r *discovResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			r.log.Errorf("[resolver] Failed to watch discov endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discovResolver) update(ins []*registry.ServiceInstance) {
	var (
		endpoints = make(map[string]struct{})
		filtered  = make([]*registry.ServiceInstance, 0, len(ins))
	)
	for _, in := range ins {
		ept, err := endpoint.ParseEndpoint(in.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		if err != nil {
			r.log.Errorf("[resolver] Failed to parse discov endpoint: %v", err)
			continue
		}
		if ept == "" {
			continue
		}
		// filter redundant endpoints
		if _, ok := endpoints[ept]; ok {
			continue
		}
		filtered = append(filtered, in)
	}

	addrs := make([]resolver.Address, 0, len(filtered))
	for _, in := range filtered {
		ept, _ := endpoint.ParseEndpoint(in.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		endpoints[ept] = struct{}{}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata).WithValue("rawServiceInstance", in),
			Addr:       ept,
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		r.log.Warnf("[resolver] Zero endpoint found,refused to write, instances: %v", ins)
		return
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		r.log.Errorf("[resolver] failed to update state: %s", err)
	}
	if r.debugLog {
		b, _ := json.Marshal(filtered)
		r.log.Infof("[resolver] update instances: %s", b)
	}
}

func (r *discovResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		r.log.Errorf("[resolver] failed to watch top: %s", err)
	}
}

func (r *discovResolver) ResolveNow(_ resolver.ResolveNowOptions) {}

func parseAttributes(md map[string]string) (a *attributes.Attributes) {
	for k, v := range md {
		a = a.WithValue(k, v)
	}
	return a
}
