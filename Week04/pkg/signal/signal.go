package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var exitChan chan struct{}
var exitOnce sync.Once

func init() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		Interrupt()
	}()
}

func Interrupt() {
	exitOnce.Do(func() {
		close(exitChan)
	})
}

func WaitInterrupt() {
	<-exitChan
}
