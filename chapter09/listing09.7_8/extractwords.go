package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	quit := make(chan struct{})
	defer close(quit)
	results := extractWords(quit, downloadedPages(quit, generateURLs(quit)))
	for result := range results {
		fmt.Println(result)
	}
}

func extractWords(quit <-chan struct{}, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)

		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, page := true, ""
		for moreData {
			select {
			case page, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(page, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
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
