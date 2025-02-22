package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	rwMutex := NewReadWriteMutex()
	for range 2 {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}

	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
}

type ReadWriteMutex struct {
	// reader lock を保持している reader の数
	readersCounter int
	// 待機している writer の数
	writerWaiting int
	// writer が lock を保持しているかどうか
	writerActive bool
	cond         *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{cond: sync.NewCond(&sync.Mutex{})}
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.cond.L.Lock()

	for rw.writerWaiting > 0 || rw.writerActive {
		// writer が待機中または lock を保持しているなら、条件変数で待機
		rw.cond.Wait()
	}
	rw.readersCounter++

	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.cond.L.Lock()

	rw.writerWaiting++
	for rw.readersCounter > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	rw.writerActive = true

	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.cond.L.Lock()

	rw.readersCounter--
	if rw.readersCounter == 0 {
		rw.cond.Broadcast()
	}

	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.cond.L.Lock()

	rw.writerActive = false
	rw.cond.Broadcast()

	rw.cond.L.Unlock()
}
