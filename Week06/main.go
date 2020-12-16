package main

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type SlidingCounter struct {
	NumBuckets int          // 包含的时间段数量，比设定值大 1
	width      int64        // 分割的每个时间段的长度
	buckets    []int64      // 包含 NumBuckets 个元素，每个元素即计数值。通过更该索引循环使用 bucket 列表
	ticker     *time.Ticker // 定时器，放到成员里方便从外部 Stop
}

func NewSlidingCounter(duration time.Duration, numBuckets int) (*SlidingCounter, error) {
	if numBuckets < 1 {
		return nil, errors.New("numBuckets 至少为 1")
	}
	segment := duration / time.Duration(numBuckets)
	numBuckets += 1 // 多一个 bucket 才可以滑动
	sc := &SlidingCounter{
		NumBuckets: numBuckets,
		width:      segment.Nanoseconds(),
		buckets:    make([]int64, numBuckets),
		ticker:     time.NewTicker(segment),
	}

	go func() {
		for range sc.ticker.C {
			idx := sc.getCurrentIndex()
			// 每次将当前 bucket 的计数值清零
			atomic.StoreInt64(&sc.buckets[idx], 0)
		}
	}()
	return sc, nil
}

// getCurrentIndex 获得当前时间下，应该使用的 bucket 的索引
func (sc *SlidingCounter) getCurrentIndex() int {
	return int(time.Now().UnixNano()/sc.width) % sc.NumBuckets
}

// Add 使计数器在当前时间对应的 bucket 的数值增加 delta
func (sc *SlidingCounter) Add(delta int) {
	idx := sc.getCurrentIndex()
	atomic.AddInt64(&sc.buckets[idx], int64(delta))
}

// Count 统计时期内的计数值
func (sc *SlidingCounter) Count() int64 {
	curIdx := sc.getCurrentIndex()
	var sum int64
	for idx := range sc.buckets {
		if idx != curIdx { // 忽略当前 bucket 的计数
			sum += atomic.LoadInt64(&sc.buckets[idx])
		}
	}
	return sum
}

func (sc *SlidingCounter) Stop() {
	sc.ticker.Stop()
}

func main() {
	sc, _ := NewSlidingCounter(time.Second, 5)
	defer sc.Stop()

	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 100)
		sc.Add(10)
		fmt.Println("current value of counter is:", sc.Count())
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Println("current value of counter is:", sc.Count())
	}
}
