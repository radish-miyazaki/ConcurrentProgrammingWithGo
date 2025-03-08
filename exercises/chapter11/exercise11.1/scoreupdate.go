package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

func main() {
	players := []*Player{
		{"James", 0, sync.Mutex{}},
		{"Ann", 0, sync.Mutex{}},
		{"Paul", 0, sync.Mutex{}},
		{"Isabel", 0, sync.Mutex{}},
		{"Peter", 0, sync.Mutex{}},
		{"Jane", 0, sync.Mutex{}},
	}
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		// Randomly select a number of players to update
		n := rand.Intn(len(players)) + 1
		rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })
		sublist := make([]*Player, n)
		copy(sublist, players[:n])

		wg.Add(1)
		go func(players []*Player) {
			defer wg.Done()

			incrementScores(players, 10)
		}(sublist)
	}
	wg.Wait()
	for _, player := range players {
		fmt.Printf("Score for %s is %d\n", player.name, player.score)
	}
}

type Player struct {
	name  string
	score int
	mutex sync.Mutex
}

func incrementScores(players []*Player, increment int) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].name < players[j].name
	})

	for _, player := range players {
		player.mutex.Lock()
	}

	for _, player := range players {
		player.score += increment
	}

	for _, player := range players {
		player.mutex.Unlock()
	}
}
