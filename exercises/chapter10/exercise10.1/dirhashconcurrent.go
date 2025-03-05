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
	sha := sha256.New()

	var prev, next *sync.WaitGroup
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		next = &sync.WaitGroup{}
		next.Add(1)
		go func(filename string, prev, next *sync.WaitGroup) {
			defer next.Done()

			fPath := filepath.Join(dir, filename)
			fmt.Println("Processing", fPath)
			hashOnFile := FHash(fPath)
			if prev != nil {
				prev.Wait()
			}

			sha.Write(hashOnFile)
		}(file.Name(), prev, next)

		prev = next
	}
	next.Wait()

	fmt.Printf("%s - %x\n", dir, sha.Sum(nil))
}

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}
