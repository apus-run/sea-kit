package discovery

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/resolver"
)

func TestBuilderErrorDiscoverer(t *testing.T) {
	notifier := &Dispatcher{}
	errMessage := errors.New("error discoverer returns error")
	discoverer := erroredDiscoverer{
		err: errMessage,
	}
	r := NewBuilder(notifier, discoverer, PrintDebugLog(false))
	_, err := r.Build(resolver.Target{}, nil, resolver.BuildOptions{})
	assert.Equal(t, errMessage, err)
}

func TestBuilderGRPCResolverRoundRobin(t *testing.T) {
	backendCount := 5

	testInstances := startTestServers(t, backendCount)
	defer testInstances.cleanup()

	notifier := &Dispatcher{}
	discoverer := FixedDiscoverer{}

	tests := []struct {
		minPeers    int
		connections int // expected number of unique connections to servers
	}{
		{minPeers: 3, connections: 3},
		{minPeers: 5, connections: 3},
		// note: test cannot succeed with minPeers < connections because resolver
		// will never return more than minPeers addresses.
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			res := NewBuilder(notifier, discoverer, WithDiscoveryMinPeers(test.minPeers), PrintDebugLog(false))

			cc, err := grpc.Dial(res.Scheme()+":///round_robin", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(GRPCServiceConfig))
			require.NoError(t, err, "could not dial using resolver's scheme")
			defer cc.Close()

			testc := grpc_testing.NewTestServiceClient(cc)

			notifier.Notify(testInstances.addresses)

			// This step is necessary to ensure that connections to all min-peers are ready,
			// otherwise round-robin may loop only through already connected peers.
			makeSureConnectionsUp(t, test.connections, testc)

			assertRoundRobinCall(t, test.connections, testc)
		})
	}
}
