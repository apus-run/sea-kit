package zlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	defer logger.Close()

	logger.Info("This is an info message")
	logger.Infof("我是日志: %v, %v", String("route", "/hello"), Int64("port", 8090))
	logger.Error("This is an error message")
}

/*
func InitLogger() logger.Logger {
	// 这里我们用一个小技巧，
	// 就是直接使用 zap 本身的配置结构体来处理
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
*/
