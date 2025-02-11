package main

import (
	"log"
	"os"
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
	if len(os.Args) < 3 {
		log.Panicln("please provide search string and the name of the files to search")
		return
	}

	searchStr := os.Args[1]

	for _, filename := range os.Args[2:] {
		go grepFile(filename, searchStr)
	}

	time.Sleep(3 * time.Second)
}
