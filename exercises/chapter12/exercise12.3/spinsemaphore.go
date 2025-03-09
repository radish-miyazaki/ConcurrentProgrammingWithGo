package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type SpinSemaphore int32

func (ss *SpinSemaphore) Acquire() {
	for {
		v := atomic.LoadInt32((*int32)(ss))

		if v != 0 && atomic.CompareAndSwapInt32((*int32)(ss), v, v-1) {
			break
		}
	}
}

func (ss *SpinSemaphore) Release() {
	atomic.AddInt32((*int32)(ss), 1)
}

func NewSpinSemaphore(permits int32) *SpinSemaphore {
	ss := new(SpinSemaphore)
	atomic.StoreInt32((*int32)(ss), permits)
	return ss
}

func acquireAndWait(id int, ss *SpinSemaphore) {
	ss.Acquire()
	fmt.Println(id, "has acquired the semaphore")
	time.Sleep(2 * time.Second)
	fmt.Println(id, "releasing the semaphore")
	ss.Release()
}

func main() {
	ss := NewSpinSemaphore(2)

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		go func(id int) {
			acquireAndWait(id, ss)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
