package env

import (
	"os"
	"strconv"

	"github.com/apus-run/sea-kit/log"
)

const (
	LllidanWebsocketPrintWriteData   = "LLLIDAN_WEBSOCKET_PRINT_WRITE"
	LllidanWebsocketPrintReadData    = "LLLIDAN_WEBSOCKET_PRINT_READ"
	LllidanWebsocketKeepAliveTimeout = "LLLIDAN_WEBSOCKET_KEEPALIVE_TIMEOUT"
)

const (
	Show   = true
	Hidden = false
)

// GetWebsocketWrite returns the debug mode of the websocket conns
func GetWebsocketWrite() bool {
	a := os.Getenv(LllidanWebsocketPrintWriteData)
	if a == "" {
		log.Infof("no env variable getting from '%s' in the container", LllidanWebsocketPrintWriteData)
		return Hidden
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		log.Error(err)
		return Hidden
	}
	return t
}

// GetWebsocketRead returns the debug mode of the websocket conns
func GetWebsocketRead() bool {
	a := os.Getenv(LllidanWebsocketPrintReadData)
	if a == "" {
		log.Infof("no env variable getting from '%s' in the container", LllidanWebsocketPrintReadData)
		return Hidden
	}
	t, err := strconv.ParseBool(a)
	if err != nil {
		log.Error(err)
		return Hidden
	}
	return t
}

// GetWebsocketKeepaliveTimeout
func GetWebsocketKeepaliveTimeout(defaultValue int64) int64 {
	a := os.Getenv(LllidanWebsocketKeepAliveTimeout)
	if a == "" {
		return defaultValue
	}
	t, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		log.Errorw("GetWebsocketKeepaliveTimeout", "default", defaultValue, "err", err)
		return defaultValue
	}
	return t
}
