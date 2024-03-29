package gs

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// 等待所有关闭信号
// Wait for all shutdown signals
func WaitingForGracefulShutdown(sigs ...*TerminateSignal) {
	quit := make(chan os.Signal, 1)                                       // 创建一个接收信号的通道 (Create a channel to receive signals)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) // 注册要接收的信号 (Register the signals to receive)
	<-quit                                                                // 等待接收信号 (Wait for the signal to receive)
	signal.Stop(quit)                                                     // 停止接收信号 (Stop receiving signals)
	close(quit)                                                           // 关闭通道 (Close the channel)
	if len(sigs) > 0 {                                                    // 执行关闭动作 (Execute the shutdown action)
		wg := sync.WaitGroup{}
		wg.Add(len(sigs))
		for _, ts := range sigs {
			go ts.Close(&wg) // 执行关闭动作 (Execute the shutdown action)
		}
		wg.Wait()
	}
}
