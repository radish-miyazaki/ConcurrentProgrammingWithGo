package main

import (
	"fmt"
	"sync"
	"time"
)

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1_000_000; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1_000_000; i++ {
		mutex.Lock()
		*money -= 10
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	mutex := sync.Mutex{}
	go stingy(&money, &mutex)
	go spendy(&money, &mutex)

	time.Sleep(2 * time.Second)
	println("Money in bank account: ", money)
}
