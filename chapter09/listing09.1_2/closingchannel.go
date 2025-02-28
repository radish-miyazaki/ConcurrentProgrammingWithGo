package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan struct{}) {
	go func() {
		for range 10 {
			fmt.Println(<-numbers)
		}
		close(quit)
	}()
}

func main() {
	numbers := make(chan int)
	quit := make(chan struct{})
	printNumbers(numbers, quit)

	next := 0
	for i := 1; ; i++ {
		next += i
		select {
		case numbers <- next:
		case <-quit:
			fmt.Println("Quitting number generation")
			return
		}
	}
}
