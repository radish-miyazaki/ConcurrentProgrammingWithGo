package main

import (
	"fmt"
	"time"
)

func main() {
	msgCh := make(chan int)
	go receiver(msgCh)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgCh <- i
		time.Sleep(1 * time.Second)
	}
	close(msgCh)
	time.Sleep(3 * time.Second)
}

func receiver(messages <-chan int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
}
