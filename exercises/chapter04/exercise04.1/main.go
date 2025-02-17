package main

import (
	"fmt"
	"sync"
	"time"
)

func countdown(seconds *int, mutex *sync.Mutex) {
	mutex.Lock()
	copiedSeconds := *seconds
	mutex.Unlock()

	for copiedSeconds > 0 {
		time.Sleep(1 * time.Second)

		mutex.Lock()
		*seconds -= 1
		copiedSeconds = *seconds
		mutex.Unlock()
	}
}

func main() {
	count := 5
	mutex := sync.Mutex{}

	go countdown(&count, &mutex)

	// 公式のドキュメントだとここは保護していなかったが、ここでもデータ競合が発生し得るため、
	// 保護した方が良いと思われる
	mutex.Lock()
	copiedCount := count
	mutex.Unlock()

	for copiedCount > 0 {
		time.Sleep(500 * time.Millisecond)
		mutex.Lock()
		fmt.Println(count)
		copiedCount = count
		mutex.Unlock()
	}
}
