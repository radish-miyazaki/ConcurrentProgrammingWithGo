package main

import (
	"fmt"
	"math/rand"
)

func main() {
	resultsCh := make([]chan []int, 10)
	for i := range 10 {
		resultsCh[i] = make(chan []int)
		go func(n int) {
			// int64 だと処理に時間がかかるため int32 に
			num := rand.Int31()
			resultsCh[i] <- findFactors(int(num))
		}(i)
	}

	for i := range 10 {
		fmt.Println(<-resultsCh[i])
	}
}

func findFactors(number int) []int {
	results := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			results = append(results, i)
		}
	}
	return results
}
