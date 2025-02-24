package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	paths := []string{}

	wg.Add(1)
	go searchFile(os.Args[1], os.Args[2], &wg, &mutex, &paths)
	wg.Wait()

	mutex.Lock()
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}

	mutex.Unlock()
}

func searchFile(
	dir string,
	filename string,
	wg *sync.WaitGroup,
	mutex *sync.Mutex,
	paths *[]string,
) {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		if strings.Contains(file.Name(), filename) {
			mutex.Lock()
			*paths = append(*paths, fpath)
			mutex.Unlock()
		}

		if file.IsDir() {
			wg.Add(1)
			go searchFile(fpath, filename, wg, mutex, paths)
		}
	}
	wg.Done()
}
