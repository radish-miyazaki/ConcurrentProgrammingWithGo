package main

import "fmt"

func main() {
	quit := make(chan struct{})
	defer close(quit)

	results := generateURLs(quit)
	for result := range results {
		fmt.Println(result)
	}
}

func generateURLs(quit <-chan struct{}) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.html", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()

	return urls
}
