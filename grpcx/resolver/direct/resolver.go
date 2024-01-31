package direct

import "google.golang.org/grpc/resolver"

type directResolver struct {
	cc resolver.ClientConn
}

func newDirectResolver(cc resolver.ClientConn) resolver.Resolver {
	return &directResolver{
		cc: cc,
	}
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(_ resolver.ResolveNowOptions) {
}
