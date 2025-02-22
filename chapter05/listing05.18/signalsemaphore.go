package main

import (
	"fmt"
	"sync"
)

func main() {
	semaphore := NewSemaphore(0)
	for range 50_000 {
		go doWork(semaphore)
		fmt.Println("Waiting for work child goroutine")
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work ended")
	semaphore.Release()
}

type Semaphore struct {
	// セマフォに残っている許可数
	permits int
	// 許可数が不足している場合に待機する際に用いる条件変数
	cond *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{permits: n, cond: sync.NewCond(&sync.Mutex{})}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	if s.permits <= 0 {
		s.cond.Wait()
	}

	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}
