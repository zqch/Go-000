package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var DefaultShutdownTimeout = 5 * time.Second // TODO: 可从环境变量读取

func ListenAndServe(addr string, handler http.Handler) error {
	// TODO: 可以再传入一些个性化的配置，这里就省略了
	srv := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		// 等待中断信号，使用 Shutdown 优雅关闭
		WaitInterrupt()
		ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	// ListenAndServe always returns a non-nil error
	return srv.ListenAndServe()
}

func ServeApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	return ListenAndServe("0.0.0.0:8080", mux)
}

func ServeDebug() error {
	return ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
}
