package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan struct{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		next = make(chan struct{})

		go func(filename string, prev, next chan struct{}) {
			fPath := filepath.Join(dir, filename)
			hashOnFile := FHash(fPath)
			if prev != nil {
				// ゴルーチンが最初の反復でなければ、前のゴルーチンが完了するまで待機
				<-prev
			}

			sha.Write(hashOnFile)
			// next <- struct{}{}
			close(next)
		}(file.Name(), prev, next)

		prev = next
	}

	<-next
	fmt.Printf("%s - %x\n", dir, sha.Sum(nil))
}

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}
