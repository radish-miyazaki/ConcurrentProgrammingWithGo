package main

import (
	"fmt"
	"net/http"
	"time"
)

type WordFrequency = map[string]int

func countWords(url string, frequency WordFrequency) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("server returning error status code: " + resp.Status)
	}

	// body, _ := io.ReadAll()

	// // 空行ごとに Split

	// for _, b := range body {
	// 	// `frequency` に word が存在していたら 1 加算
	// 	// 存在していなかったら新規追加
	// }

	fmt.Println("Completed: ", url)
}

func main() {
	var frequency WordFrequency

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countWords(url, frequency)
	}

	time.Sleep(10 * time.Second)

	for k, v := range frequency {
		fmt.Printf("%s-%d", k, v)
	}
}
