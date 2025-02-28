package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const downloaders = 20

func main() {
	quit := make(chan struct{})
	defer close(quit)

	urls := generateURLs(quit)
	pages := make([]<-chan string, downloaders)
	for i := range downloaders {
		pages[i] = downloadedPages(quit, urls)
	}

	results := longestWords(quit, extractWords(quit, FanIn(quit, pages...)))
	fmt.Println("Longest words:", <-results)
}

func longestWords(quit <-chan struct{}, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)

		// 既に出現した単語を記録するためのマップ（フィルタリングで利用）
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}

		sort.Slice(uniqueWords, func(a, b int) bool {
			return len(uniqueWords[a]) > len(uniqueWords[b])
		})

		longWords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longWords
}

func FanIn[K any](quit <-chan struct{}, allChannels ...<-chan K) <-chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))

	output := make(chan K)
	for _, channel := range allChannels {
		go func(ch <-chan K) {
			defer wg.Done()
			for i := range ch {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
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
