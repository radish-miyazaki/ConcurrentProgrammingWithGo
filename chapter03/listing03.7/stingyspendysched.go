package main

import (
	"fmt"
	"runtime"
	"time"
)

func stingy(money *int) {
	for i := 0; i < 1_000_000; i++ {
		*money += 10
		runtime.Gosched()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int) {
	for i := 0; i < 1_000_000; i++ {
		*money -= 10
		runtime.Gosched()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	go stingy(&money)
	go spendy(&money)

	time.Sleep(2 * time.Second)
	println("Money in bank account: ", money)
}
