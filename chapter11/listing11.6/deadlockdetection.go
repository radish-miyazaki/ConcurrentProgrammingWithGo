package main

import (
	"fmt"
	"sync"
)

func lockBoth(lock1, lock2 *sync.Mutex, wg *sync.WaitGroup) {
	for range 10_000 {
		lock1.Lock()
		lock2.Lock()
		lock1.Unlock()
		lock2.Unlock()
	}
	wg.Done()
}

func main() {
	lockA, lockB := sync.Mutex{}, sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go lockBoth(&lockA, &lockB, &wg)
	go lockBoth(&lockB, &lockA, &wg)
	wg.Wait()

	fmt.Println("Done")
}
