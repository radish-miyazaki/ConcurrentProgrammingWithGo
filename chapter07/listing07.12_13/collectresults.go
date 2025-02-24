package main

import "fmt"

func main() {
	resultCh := make(chan []int)
	go func() {
		resultCh <- findFactors(3_419_110_721)
	}()
	fmt.Println(findFactors(4_033_836_233))
	fmt.Println(<-resultCh)
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
