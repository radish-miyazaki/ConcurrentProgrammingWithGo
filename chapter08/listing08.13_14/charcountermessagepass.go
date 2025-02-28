package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string) <-chan []int {
	result := make(chan []int)

	go func() {
		defer close(result)
		freq := make([]int, len(allLetters))

		resp, _ := http.Get(url)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			panic("Server returning error code: " + resp.Status)
		}

		body, _ := io.ReadAll(resp.Body)
		for _, b := range body {
			c := strings.ToLower(string(b))
			cIndex := strings.Index(allLetters, c)
			if cIndex >= 0 {
				freq[cIndex] += 1
			}
		}
		fmt.Println("Completed:", url)
		result <- freq
	}()

	return result
}

func main() {
	results := make([]<-chan []int, 0)
	totalFreq := make([]int, len(allLetters))

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		results = append(results, countLetters(url))
	}

	for _, c := range results {
		freqResult := <-c
		for i := range len(freqResult) {
			totalFreq[i] += freqResult[i]
		}
	}

	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, allLetters[i])
	}
}
