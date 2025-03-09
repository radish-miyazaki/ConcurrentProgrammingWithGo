package main

import (
	"sort"
	"sync"
)

type Flight struct {
	Origin, Dest string
	SeatsLeft    int
	Locker       sync.Locker
}

func NewFlight(origin, dest string) *Flight {
	return &Flight{
		Origin:    origin,
		Dest:      dest,
		SeatsLeft: 200,
		Locker:    NewSpinLock(),
	}
}

func Book(flights []*Flight, seatsToBook int) bool {
	bookable := true

	// デッドロックを回避するために、出発地と目的地に基づいてソート
	sort.Slice(flights, func(a, b int) bool {
		flightA := flights[a].Origin + flights[a].Dest
		flightB := flights[b].Origin + flights[b].Dest
		return flightA < flightB
	})

	for _, f := range flights {
		f.Locker.Lock()
	}
	// 十分な座席があるか確認
	for i := 0; i < len(flights) && bookable; i++ {
		if flights[i].SeatsLeft < seatsToBook {
			bookable = false
		}
	}
	// 座席確保
	for i := 0; i < len(flights) && bookable; i++ {
		flights[i].SeatsLeft -= seatsToBook
	}
	for _, f := range flights {
		f.Locker.Unlock()
	}

	return bookable
}
