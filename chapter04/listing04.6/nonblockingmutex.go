package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)

	for i := 2000; i <= 2200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}

	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)

		if mutex.TryLock() {
			for i, c := range allLetters {
				fmt.Printf("%c-%d ", c, frequency[i])
			}
			mutex.Unlock()
		} else {
			fmt.Println("Mutex already being used")
		}
	}
}
