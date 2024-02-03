package discov

import (
	"context"
	"net/url"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	"github.com/apus-run/sea-kit/grpcx/registry"
)

func TestWithInsecure(t *testing.T) {
	b := &discovBuilder{}
	WithInsecure(true)(b)
	if !b.insecure {
		t.Errorf("expected insecure to be true")
	}
}

func TestWithTimeout(t *testing.T) {
	o := &discovBuilder{}
	v := time.Duration(123)
	WithTimeout(v)(o)
	if !reflect.DeepEqual(v, o.timeout) {
		t.Errorf("expected %v, got %v", v, o.timeout)
	}
}

func TestPrintDebugLog(t *testing.T) {
	o := &discovBuilder{}
	PrintDebugLog(true)(o)
	if !o.debugLog {
		t.Errorf("expected PrintdebugLog true, got %v", o.debugLog)
	}
}

type mockDiscovery struct{}

func (m *mockDiscovery) GetService(_ context.Context, _ string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (m *mockDiscovery) Watch(_ context.Context, _ string) (registry.Watcher, error) {
	time.Sleep(time.Microsecond * 500)
	return &testWatch{}, nil
}

func TestBuilder_Scheme(t *testing.T) {
	b := NewBuilder(&mockDiscovery{})
	if !reflect.DeepEqual("discov", b.Scheme()) {
		t.Errorf("expected %v, got %v", "discov", b.Scheme())
	}
}

type mockConn struct{}

func (m *mockConn) UpdateState(resolver.State) error {
	return nil
}

func (m *mockConn) ReportError(error) {}

func (m *mockConn) NewAddress(_ []resolver.Address) {}

func (m *mockConn) NewServiceConfig(_ string) {}

func (m *mockConn) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return nil
}

func TestBuilder_Build(t *testing.T) {
	b := NewBuilder(&mockDiscovery{}, PrintDebugLog(true))
	_, err := b.Build(
		resolver.Target{
			URL: url.URL{
				Scheme: resolver.GetDefaultScheme(),
				Path:   "grpc://authority/endpoint",
			},
		},
		&mockConn{},
		resolver.BuildOptions{},
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	timeoutBuilder := NewBuilder(&mockDiscovery{}, WithTimeout(0))
	_, err = timeoutBuilder.Build(
		resolver.Target{
			URL: url.URL{
				Scheme: resolver.GetDefaultScheme(),
				Path:   "grpc://authority/endpoint",
			},
		},
		&mockConn{},
		resolver.BuildOptions{},
	)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}
