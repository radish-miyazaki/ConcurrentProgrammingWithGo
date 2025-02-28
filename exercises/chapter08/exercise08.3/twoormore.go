package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	players := map[int]chan string{
		0: player(),
		1: player(),
		2: player(),
		3: player(),
	}

	for len(players) > 1 {
		select {
		case move, moreData := <-players[0]:
			handlePlayer(0, moreData, move, players)
		case move, moreData := <-players[1]:
			handlePlayer(1, moreData, move, players)
		case move, moreData := <-players[2]:
			handlePlayer(2, moreData, move, players)
		case move, moreData := <-players[3]:
			handlePlayer(3, moreData, move, players)
		}
	}

	for winner := range players {
		fmt.Printf("Game finished. Winner is player %d\n", winner)
		break
	}
}

func handlePlayer(
	id int,
	moreData bool,
	move string,
	players map[int]chan string,
) {
	if moreData {
		fmt.Printf("Player %d moved %s\n", id, move)
	} else {
		delete(players, id)
		fmt.Printf("Player %d left the game. Remaining players %d\n", id, len(players))
	}
}

func player() chan string {
	output := make(chan string)
	count := rand.Intn(100)
	move := []string{"UP", "DOWN", "LEFT", "RIGHT"}
	go func() {
		defer close(output)
		for range count {
			output <- move[rand.Intn(len(move))]
			d := time.Duration(rand.Intn(200))
			time.Sleep(d * time.Millisecond)
		}
	}()
	return output
}
