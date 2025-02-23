package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})

	playersRemaining := 5
	waitingSeconds := 3
	isTimeout := false

	go timeout(cond, waitingSeconds, &isTimeout)

	for playerId := range playersRemaining {
		go playerHandler(cond, &playersRemaining, playerId, &isTimeout)
		time.Sleep(1 * time.Second)
	}
}

// 解答例と異なるが、自分の解答の方が問題に近いと思う
func timeout(cond *sync.Cond, waitingSeconds int, isTimeout *bool) {
	time.Sleep(time.Duration(waitingSeconds) * time.Second)

	cond.L.Lock()

	// タイムアウトした場合は待機中のゴルーチンをすべて再開させる
	*isTimeout = true
	cond.Broadcast()

	cond.L.Unlock()
}

func playerHandler(
	cond *sync.Cond,
	playersRemaining *int,
	playerId int,
	isTimeout *bool,
) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	*playersRemaining--

	// タイムアウトした場合はゲームに参加できない
	if *isTimeout {
		fmt.Println(playerId, ": Game canceled")
		return
	}

	// 全プレイヤーが接続した場合は、待機中のゴルーチンを再開させる
	if *playersRemaining == 0 {
		cond.Broadcast()
	}

	// タイムアウトまたは全プレイヤーが接続するまで待機
	for *playersRemaining > 0 && !*isTimeout {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	cond.L.Unlock()

	if *isTimeout {
		// タイムアウトした場合は、接続しているユーザのみでゲーム開始
		fmt.Println("Game started. Ready player", playerId)
	} else {
		fmt.Println("All players connected. Ready player", playerId)
	}
}
