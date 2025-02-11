package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func printFileContent(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", filename, err)
	}

	fmt.Println(string(content))
}

func main() {
	if len(os.Args) < 2 {
		log.Panicln("please provide the name of the files to concatenate")
		return
	}

	for _, filename := range os.Args[1:] {
		go printFileContent(filename)
	}

	time.Sleep(3 * time.Second)
}
