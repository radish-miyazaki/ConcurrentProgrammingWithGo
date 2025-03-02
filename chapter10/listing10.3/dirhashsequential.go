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
	for _, file := range files {
		if !file.IsDir() {
			fPath := filepath.Join(dir, file.Name())
			hashOnFile := FHash(fPath)
			sha.Write(hashOnFile)
		}
	}
	fmt.Printf("%s - %x\n", dir, sha.Sum(nil))
}

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)
	return sha.Sum(nil)
}
