package main

/*
 * 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，
 * 要保证能够一个退出，全部注销退出
 */

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func main() {
	// 带取消的 context
	g, ctx := errgroup.WithContext(context.Background())
	go func() {
		// 有错误(服务退出)时，广播中断
		<-ctx.Done()
		Interrupt()
	}()

	g.Go(ServeApp)
	g.Go(ServeDebug)

	g.Wait()
}
