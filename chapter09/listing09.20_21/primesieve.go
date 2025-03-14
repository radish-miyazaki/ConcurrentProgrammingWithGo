package main

import "fmt"

func main() {
	numbers := make(chan int)
	quit := make(chan struct{})
	go primeMultipleFilter(numbers, quit)
	for i := 2; i <= 10_000; i++ {
		numbers <- i
	}
	close(numbers)
	<-quit
}

func primeMultipleFilter(numbers <-chan int, quit chan struct{}) {
	var right chan int
	p := <-numbers
	fmt.Println(p)
	for n := range numbers {
		if n%p != 0 {
			if right == nil {
				right = make(chan int)
				go primeMultipleFilter(right, quit)
			}
			right <- n
		}
	}

	if right == nil {
		close(quit)
	} else {
		close(right)
	}
}
