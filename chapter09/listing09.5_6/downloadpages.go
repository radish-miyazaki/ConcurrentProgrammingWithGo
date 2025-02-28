package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	quit := make(chan struct{})
	defer close(quit)

	results := downloadedPages(quit, generateURLs(quit))
	for result := range results {
		fmt.Println(result)
	}
}

func downloadedPages(quit <-chan struct{}, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != http.StatusOK {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	return pages
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
