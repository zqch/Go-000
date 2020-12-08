package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var exitChan = make(chan struct{})
var closeOnce sync.Once

func init() {
	var ch = make(chan os.Signal, 1)
	// SIGTERM 是 supervisor 默认的 stop signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		Interrupt()
	}()
}

// WaitInterrupt 阻塞等待中断信号
func WaitInterrupt() {
	<-exitChan
}

// Interrupt 主动让程退出。首次调用时关闭 exitChan，其他 <-ExitChan 和 WaitInterrupt 的地方不再阻塞
func Interrupt() {
	closeOnce.Do(func() {
		close(exitChan)
	})
}
