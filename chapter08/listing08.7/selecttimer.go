package main

import (
	"fmt"
	"os"
	"strconv"
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

	t, _ := strconv.Atoi(os.Args[1])
	timeoutDuration := time.Duration(t) * time.Second
	fmt.Printf("Waiting for messages for %d seconds...\n", t)
	select {
	case msg := <-msgs:
		fmt.Println("Message received:", msg)
	case tNow := <-time.After(timeoutDuration):
		fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
	}
}
