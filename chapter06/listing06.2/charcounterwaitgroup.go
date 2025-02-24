package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(31)

	mutex := sync.Mutex{}
	var frequency = make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go func() {
			countLetters(url, frequency, &mutex)
			wg.Done()
		}()
	}

	wg.Wait()

	mutex.Lock()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()

	fmt.Println("Completed: ", url)
}
