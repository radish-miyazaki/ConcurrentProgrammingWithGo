package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := NewWaitGrp(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All done")
}

func doWork(id int, wg *WaitGrp) {
	fmt.Println(id, "Done working")
	wg.Done()
}

type WaitGrp struct {
	sema *Semaphore
}

func NewWaitGrp(size int) *WaitGrp {
	return &WaitGrp{sema: NewSemaphore(1 - size)}
}

func (wg *WaitGrp) Wait() {
	wg.sema.Acquire()
}

func (wg *WaitGrp) Done() {
	wg.sema.Release()
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
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}
