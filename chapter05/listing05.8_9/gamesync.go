package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	playersRemaining := 4
	for playerId := range 4 {
		go playerHandler(cond, &playersRemaining, playerId)
		time.Sleep(1 * time.Second)
	}
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	*playersRemaining--
	if *playersRemaining == 0 {
		cond.Broadcast()
	}

	for *playersRemaining > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}

	cond.L.Unlock()
	fmt.Println("All players connected. Ready player", playerId)
}
