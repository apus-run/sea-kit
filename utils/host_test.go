package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLocalIp(t *testing.T) {
	ips := []struct {
		host   string
		expect bool
	}{
		{
			"192.168.1.1:9090",
			false,
		},
		{
			"127.0.0.1:8080",
			true,
		},
		{
			"localhost:8080",
			true,
		},
	}

	for _, ip := range ips {
		if IsLocalIP(ip.host) == ip.expect {
			t.Logf("ip:host %s pass", ip.host)
		} else {
			t.Errorf("ip:host %s is not expected!", ip.host)
		}
	}
}

func TestGetAvailablePort(t *testing.T) {
	port, err := GetAvailablePort()
	assert.NoError(t, err)
	t.Log(port)
}

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	t.Log(hostname)
}

func TestGetLocalHTTPAddrPairs(t *testing.T) {
	serverAddr, requestAddr := GetLocalHTTPAddrPairs()
	t.Logf("\t serverAddr: %v\n requestAddr: %v", serverAddr, requestAddr)
	assert.NotEmpty(t, serverAddr)
	assert.NotEmpty(t, requestAddr)
}
