package prof

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/apus-run/sea-kit/utils"
)

func TestRegister(t *testing.T) {
	r := gin.Default()
	Register(r, WithPrefix(""), WithPrefix("/myServer"), WithIOWaitTime())

	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("listen and serve error: " + err.Error())
		}
	}()
	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get(requestAddr + "/myServer")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPprof(t *testing.T) {
	mux := New()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			for i := 0; i < 10000; i++ {
				log.Println("current index: ", i)
				time.Sleep(200 * time.Millisecond)
			}
		}()

		w.Write([]byte("hello"))
	})

	Run(mux, 8080)

	ch := make(chan os.Signal, 1)
	log.Println("wait exit signal...")
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal support syscall.SIGUSR2
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch
	log.Println("exit signal: ", sig.String())

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()

	log.Println("shutting down")
}
