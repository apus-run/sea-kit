package env

import (
	"os"
	"strconv"

	"github.com/apus-run/sea-kit/log"
)

const (
	LllidanGatewayPProfDebug = "LLLIDAN_GATEWAY_PPROF_DEBUG"
)

// GetGatewayPProfDebug determines that whether the gateway will open pprof routers
func GetGatewayPProfDebug() bool {
	a := os.Getenv(LllidanGatewayPProfDebug)
	if a == "" {
		return false
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		log.Error(err)
		return false
	}
	return t
}
