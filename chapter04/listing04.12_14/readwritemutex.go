package main

import "sync"

type ReadWriteMutex struct {
	// クリティカルセクション内にあるリーダゴルーチンの数
	readerCounter int
	// リーダゴルーチンの数を保持するためのロック
	readersLock sync.Mutex
	// グローバルロック
	globalLock sync.Mutex
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.readersLock.Lock()

	rw.readerCounter++
	if rw.readerCounter == 1 {
		// リーダゴルーチンがクリティカルセクション内に始めて入った場合、`globalLock` をロック
		rw.globalLock.Lock()
	}

	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.readersLock.Lock()

	rw.readerCounter--
	if rw.readerCounter == 0 {
		rw.globalLock.Unlock()
	}

	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}
