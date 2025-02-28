package main

import (
	"fmt"
	"time"
)

func writeEvery(msg string, seconds time.Duration) <-chan string {
	msgs := make(chan string)
	go func() {
		for {
			time.Sleep(seconds)
			msgs <- msg
		}
	}()
	return msgs
}

func main() {
	msgsFromA := writeEvery("Tick", 1*time.Second)
	msgsFromB := writeEvery("Tack", 3*time.Second)

	for {
		select {
		case msg1 := <-msgsFromA:
			fmt.Println(msg1)
		case msg2 := <-msgsFromB:
			fmt.Println(msg2)
		}
	}
}
