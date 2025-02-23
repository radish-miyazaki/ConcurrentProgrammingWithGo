package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ws := NewWeightedSemaphore(3)
	fmt.Printf("Parent thread created semaphore with %d permits\n", 3)

	ws.Acquire(2)
	fmt.Printf("Parent thread acquired %d semaphores\n", 2)

	go func() {
		ws.Acquire(2)
		fmt.Printf("Child thread acquired %d semaphores\n", 2)

		ws.Release(2)
		fmt.Printf("Child thread released %d semaphores\n", 2)
	}()

	time.Sleep(3 * time.Second)

	fmt.Printf("Parent thread releasing %d semaphores\n", 2)
	ws.Release(2)

	time.Sleep(1 * time.Second)
}

type WeightedSemaphore struct {
	// セマフォに残っている許可数
	permits int
	// 許可数が不足している場合に待機する際に用いる条件変数
	cond *sync.Cond
}

func NewWeightedSemaphore(n int) *WeightedSemaphore {
	return &WeightedSemaphore{permits: n, cond: sync.NewCond(&sync.Mutex{})}
}

func (ws *WeightedSemaphore) Acquire(permits int) {
	ws.cond.L.Lock()
	for ws.permits < permits {
		ws.cond.Wait()
	}
	ws.permits -= permits
	ws.cond.L.Unlock()
}

func (ws *WeightedSemaphore) Release(permits int) {
	ws.cond.L.Lock()
	ws.permits += permits
	ws.cond.Signal()
	ws.cond.L.Unlock()
}
