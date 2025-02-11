package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func grepFile(filename string, searchStr string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", filename, err)
	}

	if strings.Contains(string(content), searchStr) {
		log.Printf("file %s contains search string %s\n", filename, searchStr)
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
		if file.IsDir() {
			continue
		}

		go grepFile(filepath.Join(dirPath, file.Name()), searchStr)
	}

	time.Sleep(3 * time.Second)
}
