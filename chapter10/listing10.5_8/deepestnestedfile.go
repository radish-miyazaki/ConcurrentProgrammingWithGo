package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	dir := os.Args[1]
	partialResults := make(chan CodeDepth)
	wg := &sync.WaitGroup{}

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		forkIfNeeded(path, info, wg, partialResults)
		return nil
	})
	finalResult := joinResults(partialResults)

	wg.Wait()
	close(partialResults)

	result := <-finalResult
	fmt.Printf("%s has the deepest nested code block of %d\n", result.filename, result.level)
}

type CodeDepth struct {
	filename string
	level    int
}

func deepestNestedBlock(filename string) CodeDepth {
	code, _ := os.ReadFile(filename)
	max := 0
	level := 0

	for _, c := range code {
		if c == '{' {
			level++
			max = int(math.Max(float64(max), float64(level)))
		} else if c == '}' {
			level--
		}
	}
	return CodeDepth{filename: filename, level: max}
}

func forkIfNeeded(
	path string,
	info os.FileInfo,
	wg *sync.WaitGroup,
	results chan CodeDepth,
) {
	if !info.IsDir() && strings.HasSuffix(path, ".go") {
		wg.Add(1)

		go func() {
			results <- deepestNestedBlock(path)
			wg.Done()
		}()
	}
}

func joinResults(partialResults chan CodeDepth) chan CodeDepth {
	finalResult := make(chan CodeDepth)

	go func() {
		max := CodeDepth{"", 0}

		for pr := range partialResults {
			if pr.level > max.level {
				max = pr
			}
		}

		finalResult <- max
	}()

	return finalResult
}
