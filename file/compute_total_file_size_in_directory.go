package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

func walkDir(dir string, sizeChan chan<- int64) {
	defer wg.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, sizeChan)
			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			//fmt.Printf("Visited: %s File name: %s Size: %d bytes\n", path, f.Name(), f.Size())
			sizeChan <- f.Size()
		}
		return nil
	}

	filepath.Walk(dir, visit)
}

func main() {
	if len(os.Args) < 2 {
		println("Usage: compute_total_file_size_in_directory [path]")
		os.Exit(1)
	}
	root := os.Args[1]
	var sizeChan = make(chan int64)

	wg.Add(1)
	var sum int64 = 0
	go func() {
		for size := range sizeChan {
			sum += size
		}
	}()

	walkDir(root, sizeChan)
	wg.Wait()
	fmt.Println(sum)
}
