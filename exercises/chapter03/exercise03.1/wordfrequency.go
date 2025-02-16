package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type WordFrequency = map[string]int

func countWords(url string, frequency WordFrequency) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	words := regexp.MustCompile("[a-zA-Z]+").FindAllString(string(body), -1)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		frequency[lowerWord] += 1
	}

	fmt.Println("Completed: ", url)
}

func main() {
	var frequency WordFrequency = map[string]int{}

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		countWords(url, frequency)
	}

	time.Sleep(10 * time.Second)

	for k, v := range frequency {
		fmt.Printf("%s -> %d\n", k, v)
	}
}
