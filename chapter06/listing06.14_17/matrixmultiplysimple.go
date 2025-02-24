package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const matrixSize = 3

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int

	barrier := NewBarrier(matrixSize + 1)
	for row := range matrixSize {
		go rowMultiply(&matrixA, &matrixB, &result, row, barrier)
	}

	for range 3 {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)

		// ゴルーチンが計算できるようにバリアを開放
		barrier.Wait()

		// ゴルーチンが計算を完了するまで待機
		barrier.Wait()

		for i := range matrixSize {
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println("")
	}
}

func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		for col := range matrixSize {
			sum := 0
			for i := range matrixSize {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
		}
	}
}

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		for col := range matrixSize {
			// -5 ~ 4 のランダムな値を割り当てる
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(
	matrixA,
	matrixB,
	result *[matrixSize][matrixSize]int,
	row int,
	barrier *Barrier,
) {
	for {
		// main ゴルーチンが行列を読み込むまで待機
		barrier.Wait()

		for col := range matrixSize {
			sum := 0
			for i := range matrixSize {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
		}

		// すべての行の計算が終わるまで待機
		barrier.Wait()
	}
}

type Barrier struct {
	size      int
	waitCount int
	cond      *sync.Cond
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size, 0, sync.NewCond(&sync.Mutex{}),
	}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount++
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}

	b.cond.L.Unlock()
}
