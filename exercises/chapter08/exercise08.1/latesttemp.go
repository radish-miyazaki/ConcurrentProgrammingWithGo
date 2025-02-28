package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	tempCh := generateTemp()
	display := make(chan int)
	outputTemp(display)
	t := <-tempCh

	for {
		select {
		case t = <-tempCh:
		case display <- t:
		}
	}
}

func generateTemp() chan int {
	output := make(chan int)
	go func() {
		temp := 50 // 華氏
		for {
			output <- temp
			temp += rand.Intn(3) - 1
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func outputTemp(input chan int) {
	go func() {
		for {
			fmt.Println("Current temp:", <-input)
			time.Sleep(2 * time.Second)
		}
	}()
}
