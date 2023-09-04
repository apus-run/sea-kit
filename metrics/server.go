package metrics

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/apus-run/sea-kit/log"
	"github.com/apus-run/sea-kit/tls"
)

type Config struct {
	ServerAddress   string
	ServerTLSConfig *tls.Config
	EnableProfiling bool
}

const (
	defaultBindAddress = ":8080"
	metricsPath        = "/metrics"
)

func Serve(ctx context.Context, config Config) {
	if config.ServerAddress == "" {
		config.ServerAddress = defaultBindAddress
	}
	if config.ServerAddress == "0" {
		return
	}

	log.Infof("metrics server is starting to listen at %s", config.ServerAddress)
	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatalf("error creating the metrics listener: %v", err)
	}

	handler := promhttp.HandlerFor(Registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.HTTPErrorOnError,
	})
	mux := http.NewServeMux()
	mux.Handle(metricsPath, handler)

	if config.EnableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	server := http.Server{
		Handler: mux,
	}

	go func() {
		log.Infof("starting metrics server path %s", metricsPath)
		var err error
		if config.ServerTLSConfig.Cert != "" && config.ServerTLSConfig.Key != "" {
			err = server.ServeTLS(listener, config.ServerTLSConfig.Cert, config.ServerTLSConfig.Key)
		} else {
			err = server.Serve(listener)
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting the metrics server: %v", err)
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down the metrics server: %v", err)
	}
}
