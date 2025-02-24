package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	msgCh := make(chan int, 3)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go receiver(msgCh, &wg)
	for i := 1; i <= 6; i++ {
		size := len(msgCh)
		fmt.Printf("%s Sending: %d. Buffer size: %d\n", time.Now().Format("15:04:05"), i, size)
		msgCh <- i
		fmt.Printf("%s Sent: %d. Buffer size: %d\n", time.Now().Format("15:04:05"), i, size)
	}
	msgCh <- -1
	wg.Wait()
}

func receiver(messages chan int, wg *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = <-messages
		fmt.Println("Received:", msg)
	}
	wg.Done()
}
