package atexit

import (
	"fmt"
	"sync"

	signal "github.com/zqch/Go-000/Week04/pkg/signal"
)

var exitFuncs []func()
var exitMu sync.Mutex

func init() {
	exitFuncs = make([]func(), 0)
	go func() {
		signal.WaitInterrupt()
	}()
}

func doExitFuncs() {
	for _, f := range exitFuncs {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("atexit error: %s\n", err)
			}
		}()
		f()
	}
}

func Register(f func()) {
	exitMu.Lock()
	exitFuncs = append(exitFuncs, f)
	exitMu.Unlock()
}
