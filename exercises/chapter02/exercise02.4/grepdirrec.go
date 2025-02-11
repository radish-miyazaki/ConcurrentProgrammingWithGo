package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func grepPath(basePath string, entry os.DirEntry, searchStr string) {
	fullPath := filepath.Join(basePath, entry.Name())
	if entry.IsDir() {
		files, err := os.ReadDir(fullPath)
		if err != nil {
			log.Fatalf("failed to read directory %s: %v", fullPath, err)
			return
		}

		for _, file := range files {
			go grepPath(fullPath, file, searchStr)
		}

		return
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", fullPath, err)
	}

	if strings.Contains(string(content), searchStr) {
		log.Printf("file %s contains search string %s\n", fullPath, searchStr)
	}
}

func main() {
	if len(os.Args) == 2 {
		log.Panicln("please provide search string and the name of one directory to search")
		return
	}

	searchStr := os.Args[1]
	dirPath := os.Args[2]

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("failed to read directory %s: %v", os.Args[2], err)
	}

	for _, file := range files {
		go grepPath(dirPath, file, searchStr)
	}

	time.Sleep(3 * time.Second)
}
