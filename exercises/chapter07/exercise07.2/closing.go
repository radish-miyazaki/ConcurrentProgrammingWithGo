package main

import (
	"fmt"
)

func main() {
	strCh := make(chan string, 1)
	sliceCh := make(chan []string, 1)
	close(strCh)
	close(sliceCh)
	fmt.Println("string close channel:", <-strCh)
	fmt.Println("slice close channel:", <-sliceCh)
}
