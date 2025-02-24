package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	// バリアで管理するゴルーチン数
	size int
	// 待機中のゴルーチン数
	waitCount int
	// 条件変数
	cond *sync.Cond
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size, 0, sync.NewCond(&sync.Mutex{}),
	}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount++
	if b.waitCount == b.size {
		// すべてのゴルーチンがバリアに達した場合、全ゴルーチンを再開させる
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}

	b.cond.L.Unlock()
}

func main() {
	barrier := NewBarrier(2)
	go workAndWait("Red", 4, barrier)
	go workAndWait("Blue", 10, barrier)
	time.Sleep(100 * time.Second)
}

func workAndWait(name string, timeToWait int, barrie *Barrier) {
	start := time.Now()

	for {
		fmt.Println(time.Since(start), name, "is running")

		// 指定された秒数だけ作業している状態をシミュレーション
		time.Sleep(time.Duration(timeToWait) * time.Second)

		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrie.Wait()
	}
}
