package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	filesCh := make(chan string)
	dirsCh := make(chan string)

	go handleFiles(filesCh, dirsCh)
	go handleDirectories(dirsCh, filesCh)

	dirsCh <- os.Args[1]

	time.Sleep(60 * time.Second)
}

func handleDirectories(dirs <-chan string, files chan<- string) {
	for fullpath := range dirs {
		fmt.Println("Reading all files from", fullpath)
		filesInDir, _ := os.ReadDir(fullpath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullpath)

		for _, file := range filesInDir {
			go func(fp string) {
				files <- fp
			}(filepath.Join(fullpath, file.Name()))
		}
	}
}

func handleFiles(files <-chan string, dirs chan<- string) {
	for path := range files {
		file, _ := os.Open(path)
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirs <- path
			continue
		}

		fmt.Printf("File %s, size: %.2fKB, last modified: %s\n",
			fileInfo.Name(), float64(fileInfo.Size())/1024.0,
			fileInfo.ModTime().Format(time.ANSIC))
	}
}
