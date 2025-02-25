package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Channel[M any] struct {
	cond        *sync.Cond
	maxCapacity int
	buffer      *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		maxCapacity: capacity,
		buffer:      list.New(),
		cond:        sync.NewCond(&sync.Mutex{}),
	}
}

func (c *Channel[M]) Send(message M) {
	c.cond.L.Lock()
	for c.buffer.Len() == c.maxCapacity {
		c.cond.Wait()
	}
	c.buffer.PushBack(message)
	c.cond.Signal()
	c.cond.L.Unlock()
}

func (c *Channel[M]) Receive() M {
	c.cond.L.Lock()

	// バッファなしチャネルであっても動作させるために、最初に maxCapacity をインクリメントする
	c.maxCapacity++
	c.cond.Signal()

	for c.buffer.Len() == 0 {
		c.cond.Wait()
	}

	c.maxCapacity-- // 1 加算したので、1 減算する
	// mutex で競合状態から保護しながら、バッファからメッセージを 1 つ取り出す
	v := c.buffer.Remove(c.buffer.Front()).(M)

	c.cond.L.Unlock()
	return v
}

func receiver(ch *Channel[int], wg *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = ch.Receive()
		fmt.Println("Received:", msg)
	}
	wg.Done()
}

func main() {
	// ch := NewChannel[int](0)
	ch := NewChannel[int](2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go receiver(ch, &wg)
	for i := 1; i <= 6; i++ {
		fmt.Println("Sending:", i)
		ch.Send(i)
	}
	ch.Send(-1)
	wg.Wait()
}
