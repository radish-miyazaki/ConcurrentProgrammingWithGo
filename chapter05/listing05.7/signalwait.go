package main

import (
	"fmt"
	"sync"
)

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	for range 50000 {
		go doWork(cond)
		fmt.Println("Waiting for work child goroutine ")
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}
	cond.L.Unlock()
}
