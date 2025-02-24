package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const matrixSize = 3

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int

	for range 3 {
		wg := sync.WaitGroup{}
		wg.Add(matrixSize)

		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)

		for row := range matrixSize {
			go rowMultiply(&matrixA, &matrixB, &result, row, &wg)
		}

		// ゴルーチンが計算を完了するまで待機
		wg.Wait()

		for i := range matrixSize {
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println("")
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
	wg *sync.WaitGroup,
) {
	for col := range matrixSize {
		sum := 0
		for i := range matrixSize {
			sum += matrixA[row][i] * matrixB[i][col]
		}
		result[row][col] = sum
	}

	// すべての行の計算が終わるまで待機
	wg.Done()
}
