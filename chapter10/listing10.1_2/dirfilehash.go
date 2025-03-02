package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}

	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)

			go func(filename string) {
				fPath := filepath.Join(dir, filename)
				hash := FHash(fPath)
				fmt.Printf("%s - %x\n", filename, hash)
				wg.Done()
			}(file.Name())
		}
	}

	wg.Wait()
}

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}
