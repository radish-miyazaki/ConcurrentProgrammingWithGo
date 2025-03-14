package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

func sender(messages chan<- int) {
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		messages <- i
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgCh := make(chan int)
	go receiver(msgCh)
	go sender(msgCh)
	time.Sleep(5 * time.Second)
}
