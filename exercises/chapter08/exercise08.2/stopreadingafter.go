package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	numCh := generateNumber()
	timeout := time.After(5 * time.Second)

	for {
		select {
		case num := <-numCh:
			fmt.Println("generateNumber returns", num)
		case <-timeout:
			fmt.Println("Timeout")
			return
		}
	}
}

func generateNumber() chan int {
	output := make(chan int)
	go func() {
		for {
			output <- rand.Intn(10)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}
