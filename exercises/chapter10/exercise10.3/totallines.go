package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// func main() {
// 	const pagesToDownload = 30
// 	totalLines := 0
//
// 	for i := 1000; i < 1000+pagesToDownload; i++ {
// 		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
// 		fmt.Println("Downloading", url)
// 		resp, _ := http.Get(url)
// 		if resp.StatusCode != http.StatusOK {
// 			panic("Server's error:" + resp.Status)
// 		}
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		totalLines += strings.Count(string(bodyBytes), "\n")
// 		resp.Body.Close()
// 	}
// 	fmt.Println("Total lines:", totalLines)
// }

// apply fork/join pattern
func main() {
	const pagesToDownload = 30

	countCh := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 1000; i < 1000+pagesToDownload; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			fmt.Println("Downloading", url)

			resp, _ := http.Get(url)
			if resp.StatusCode != http.StatusOK {
				panic("Server's error:" + resp.Status)
			}

			bodyBytes, _ := io.ReadAll(resp.Body)
			countCh <- strings.Count(string(bodyBytes), "\n")
			resp.Body.Close()
		}()
	}

	totalLinesCh := make(chan int)
	go func() {
		totalLines := 0

		for c := range countCh {
			totalLines += c
		}

		totalLinesCh <- totalLines
	}()

	wg.Wait()
	close(countCh)

	fmt.Println("Total lines:", <-totalLinesCh)
}
