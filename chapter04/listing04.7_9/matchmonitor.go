package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// matchRecorder は定期的な試合イベントをシミュレートする試合記録用の関数
func matchRecorder(matchEvents *[]string, mutex *sync.Mutex) {
	for i := 0; ; i++ {
		mutex.Lock()
		*matchEvents = append(*matchEvents, "Match event "+strconv.Itoa(i))
		mutex.Unlock()

		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, mutex *sync.Mutex, st time.Time) {
	// 同じユーザが複数のリクエストを行うことをシミュレート
	for i := 0; i < 100; i++ {
		mutex.Lock()
		// 試合イベントのスライス全体をコピーし、クライアントへのレスポンスをシミュレート
		allEvents := copyAllEvents(mEvents)
		mutex.Unlock()

		// 開始からの所要時間を計算
		timeTaken := time.Since(st)
		fmt.Println(len(allEvents), " events copied in ", timeTaken)
	}
}

func copyAllEvents(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}

	return allEvents
}

func main() {
	mutex := sync.Mutex{}
	var matchEvents = make([]string, 0, 10_000)

	// 進行中の試合をシミュレートする
	for j := 0; j < 10_000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}

	// 試合記録のゴルーチンを起動
	go matchRecorder(&matchEvents, &mutex)

	// クライアントのリクエストをシミュレート
	start := time.Now()
	for j := 0; j < 5_000; j++ {
		go clientHandler(&matchEvents, &mutex, start)
	}

	time.Sleep(100 * time.Second)
}
