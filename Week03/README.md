# Concurrency
## 代码作业
> 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出

见 [main.go](main.go)

## 使用 goroutine 的姿势
golang 使用 goroutine 来并发，使用 goroutine 时需要注意:

- Leave concurrency to the caller: 把并发交给调用者
- Never start a goroutine without knowing when it will stop: 需要管理起 goroutine 的生命周期，知道何时退出

## 基本并发控制
并发条件下，通常要对临界区的代码进行控制，加以保护。sync 包提供了许多同步原语。

- Mutex: 互斥锁。临界区只能由一个 goroutine 持有
- RWMutex: 读写锁。可以有多个 goroutine 同时读，但只能有一个 goroutine 写，有人写的时候其他人不能读
- WaitGroup: 阻塞等待多个操作完成
- Cond: 条件变量。不咋用，用起来也麻烦
- Once: 保证只执行一次
- Pool: 池。保存临时对象，减小创建开销，但随时可能会被回收
- Atomic: 原子操作。由硬件实现，是实现锁的底层依赖，更轻量

### 黑魔法
使用 atomic 对 map 复制是原子的，Copy-On-Write 的思路可以利用这个特点，来达到对 map 数据的无锁并发访问

单个机器字(machine word)的写入是原子的，超过单个机器字的写入不是原子的。对 interface 的赋值，看起来是原子的，实际却不是。interface 内部有两个指针，分别表示 Type 和 Data，是两个机器字

## channel
"Don't communicate by sharing memory, share memory by communicating"

区别于传统的共享锁来控制并发的方式，golang 推荐使用 channel 在 goroutine 之间通信来共享数据。

channel 的缓冲区可以降低收发延迟，但并不能提高吞吐量

### channel 并发姿势
**Generator**
```go
func work() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		// ...
		ch <- "result"
	}()
	return ch
}
```

**Pipeline**
```go
for v := <-input {
	output <- v
}
```

**Timing out**
```go
select {
case s := <-c:
	fmt.Println(s)
case <-time.After(1 * time.Second):
	fmt.Println("timeout")
	return
}
```

**Fan-in**
```go
func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1: c <- s
			case s := <-input2: c <- s
			}
		}
	}()
	return c
}
```

**Or-Done**
```go
func First(query1, query2 <-chan Result) Result {
	c := make(chan Result)
	go func() {
		select {
		case ret := <-query1: c <-ret
		case ret := <-query2: c <-ret
		}
	}()
	return <-c
}
```

**Fan-out**
```go
func fanOut(input <-chan string, outputs []chan string) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(outputs[i])
			}
		}()

		for v := range input {
			for i := 0; i < len(outputs); i++ {
				outputs[i] <- v
			}
		}
	}()
}
```

## 内存模型
内存模型描述的是并发环境中，多 goroutine 读相同变量时变量的可见性条件。在并发环境中，由于 CPU 指令重排、多级 Cache、编译器指令重排的存在，机器真正的执行顺序可能和代码的表达不一致。

### happens-before
如果写操作 w 先于读操作 r 发生，就说 w **happens-before** r，或者说 r **happens-after** w。如果 r 既不 happens before w，也不 happens after w，就说 r 和 w 是并发的

### 顺序
- 普通代码:
	在单个 goroutine 中，读写的顺序就和代码表现的顺序一样

- 初始化:
	如果一个 package p 依赖 package q，则 q 的 `init` 函数要 happens before 任何一个 p 的 `init` 函数。`main.main` 的启动 happens after 所有的 `init` 函数执行完

- goroutine 创建:
	go 语句 happens before 所创建的 goroutine 的执行

- goroutine 销毁:
	一个 goroutine 的退出不能保证 happens before 任何事件，需要显式同步机制来保证

- Channel 通信:
	channel 的关闭 happens before 收到一个由于关闭而带来的零值
	对于不带缓冲的 channel, 接收行为 happens before 对应的发送行为完成
	对于带缓冲的 channel, 发送行为 happens before 对应的接收行为完成
	对于容量为 C 的 channel, 第 k 个接收行为 happens before 第 k+C 个发送行为完成

- 锁:
	如果 n 小于 m，第 n 次 `Unlock()` 调用 happens before 第 m 次 `Lock()` 的返回
	如果某个 `RLock` 的调用 happens after 第 n 次调用 `Unlock`，则它对应的 `RUnlock` 调用 happens before 第 n+1 次 `Lock` 调用

- Once:
	`once.Do(f)` 里面对 `f()` 的调用 happens before 任意 `once.Do(f)` 的返回

## 并发的高阶姿势
**errgroup**

用于将多个小任务并发执行，底层基于 WaitGroup 实现

```go
var g errgroup.Group

// 启动三个子任务，其中一个执行失败
g.Go(func() error {
	fmt.Println("exec #1")
	return nil
})

g.Go(func() error {
	fmt.Println("exec #2")
	return errors.New("failed to exec #2")
})

g.Go(func() error {
	fmt.Println("exec #3")
	return nil
})

// 等待三个任务都完成，拿到错误结果
if err := g.Wait(); err == nil {
	fmt.Println("Successfully exec all")
} else {
	fmt.Println("failed:", err)
}
```

缺点:
- 没有限制 goroutine 的数量
- 一个子任务失败，不会取消正在执行的其他任务
- 没有对 panic 进行 recover

**SingleFlight**

用于在多个 goroutine 同时调用同一个函数时，合并成一个调用，将这一次调用的结果返回给同时调用的其他 goroutine。可以用来解决缓存击穿
