package main

import (
	"fmt"
	"time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
	msgs := make(chan string)
	go func() {
		time.Sleep(seconds)
		msgs <- "Hello"
	}()
	return msgs
}

func main() {
	msgs := sendMsgAfter(3 * time.Second)
	for {
		select {
		case msg := <-msgs:
			fmt.Println("Message received:", msg)
			return
		default:
			fmt.Println("No messages waiting")
			time.Sleep(1 * time.Second)
		}
	}
}
