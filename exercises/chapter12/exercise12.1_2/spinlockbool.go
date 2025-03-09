package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type SpinLock atomic.Bool

func (s *SpinLock) Lock() {
	for !(*atomic.Bool)(s).CompareAndSwap(false, true) {
		// Go scheduler を呼び出して、他のゴルーチンに実行を譲る（yield）
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	(*atomic.Bool)(s).Store(false)
}

func (s *SpinLock) TryLock() bool {
	return (*atomic.Bool)(s).CompareAndSwap(false, true)
}

func NewSpinLock() *SpinLock {
	var lock SpinLock
	return &lock
}

func main() {
	spinLock := NewSpinLock()
	spinLock.Lock()
	fmt.Println("The should be false:", spinLock.TryLock())
	spinLock.Unlock()
	fmt.Println("The should be true:", spinLock.TryLock())
	spinLock.Unlock()
}
