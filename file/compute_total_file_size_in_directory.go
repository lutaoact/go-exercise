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

	Walk(dir, visit)
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	//sort.Strings(names) //拿掉没有用的排序逻辑
	return names, nil
}

func walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	names, err := readDirNames(path)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func Walk(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
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
