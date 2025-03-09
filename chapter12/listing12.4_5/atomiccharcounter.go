package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	var frequency = make([]int32, 26)

	wg := sync.WaitGroup{}
	wg.Add(31)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go func() { countLetters(url, frequency); wg.Done() }()
	}
	wg.Wait()

	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, atomic.LoadInt32(&frequency[i]))
	}
}

func countLetters(url string, frequency []int32) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Server returning error code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIdx := strings.Index(allLetters, c)
		if cIdx >= 0 {
			atomic.AddInt32(&frequency[cIdx], 1)
		}
	}

	fmt.Println("Completed:", url)
}
