package main

import (
	"container/list"
	"sync"
)

type Channel[M any] struct {
	// バッファがいっぱいになったときに、送信側を待たせる容量セマフォ
	capacitySema *Semaphore
	// バッファが空になったときに、受信側を待たせるバッファサイズセマフォ
	sizeSema *Semaphore
	// 共有リストデータ構造を保護するミューテックス
	mutex *sync.Mutex
	// キューデータ構造として用いるリンクリスト
	buffer *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		capacitySema: NewSemaphore(0),
		sizeSema:     NewSemaphore(capacity),
		mutex:        &sync.Mutex{},
		buffer:       list.New(),
	}
}

func (c *Channel[M]) Send(message M) {
	// 容量セマフォから 1 つの許可を要求
	c.capacitySema.Acquire()

	// mutex で競合状態から保護しながら、バッファにメッセージを 1 つ追加
	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	// バッファサイズセマフォから 1 つの許可を解放
	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	// 容量セマフォから 1 つの許可を解放
	// バッファなしチャネルであっても動作させるために、最初に Release を呼び出す
	c.capacitySema.Release()

	// バッファサイズセマフォから 1 つの許可を要求
	c.sizeSema.Acquire()

	// mutex で競合状態から保護しながら、バッファからメッセージを 1 つ取り出す
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()

	return v
}

/*
	Semaphore
*/

type Semaphore struct {
	// セマフォに残っている許可数
	permits int
	// 許可数が不足している場合に待機する際に用いる条件変数
	cond *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{permits: n, cond: sync.NewCond(&sync.Mutex{})}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}
