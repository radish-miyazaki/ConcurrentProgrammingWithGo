package main

import (
	"fmt"
)

// 書籍で期待する main 関数
//
// 起こり得る問題点
// - TakeUntil 内で quit チャネルを閉じることで、GenerateSquares だけでなく、Print や Drain も同じキャンセルシグナルを受け取る
// - 結果、Print や Drain の処理が中途半端にキャンセルされる可能性がある
// func main() {
// 	quitChannel := make(chan struct{})
//
// 	Drain(quitChannel,
// 		Print(quitChannel,
// 			TakeUntil(func(s int) bool { return s <= 1_000_000 }, quitChannel,
// 				GenerateSquares(quitChannel))))
//
// 	<-quitChannel
// }

// 改良版（listing09.18_19 を参考）
func main() {
	squaresQuitCh := make(chan struct{})
	quitCh := make(chan struct{})
	defer close(quitCh)

	Drain(quitCh,
		Print(quitCh,
			TakeUntil(func(s int) bool { return s <= 1_000_000 }, squaresQuitCh,
				GenerateSquares(squaresQuitCh))))
}

func GenerateSquares(quit <-chan struct{}) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 1; ; i++ {
			select {
			case ch <- i * i:
			case <-quit:
				return
			}
		}
	}()

	return ch
}

func TakeUntil[K any](
	f func(K) bool,
	quit chan struct{},
	input <-chan K,
) <-chan K {
	ch := make(chan K)

	go func() {
		// この関数が終了した時点で、quit チャネルを閉じることで
		defer close(ch)
		defer close(quit)

		moreData := true
		fVal := true
		var v K

		for fVal && moreData {
			select {
			case v, moreData = <-input:
				fVal = f(v)

				// 次の値が存在しないか、条件を満たさない場合はメッセージを送信しない
				if !moreData || !fVal {
					return
				}

				ch <- v
			case <-quit:
				return
			}
		}
	}()

	return ch
}

func Print[T any](quit <-chan struct{}, input <-chan T) <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for {
			select {
			case v, moreData := <-input:
				if !moreData {
					return
				}

				fmt.Println(v)
				ch <- v
			case <-quit:
				return
			}
		}
	}()

	return ch
}

func Drain[T any](quit <-chan struct{}, input <-chan T) {
	for {
		select {
		case _, moreData := <-input:
			if !moreData {
				return
			}
		case <-quit:
			return
		}
	}
}
