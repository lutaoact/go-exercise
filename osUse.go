package main

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

func main() {
	List("/tmp/tmp.txt")
}

func List(filepath string) ([]string, error) {
	dir, err := os.Open(filepath)
	fmt.Println("open", dir, err)
	if err != nil {
		fmt.Println(syscall.ENOENT == err)
		if os.IsNotExist(err) {
			fmt.Println("not exist", dir, err)
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(0)
	fmt.Println("Readdirnames", fileNames, err)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames {
		keys = append(keys, path.Join(filepath, fileName))
	}

	fmt.Println(keys)
	return keys, nil
}

func Stat() {
	fi, err := os.Stat("/tmp")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("fi = %+v\n", fi)
	fmt.Printf("fi.IsDir() = %+v\n", fi.IsDir())

	fi, err = os.Stat("/tmp/tmp.txt")
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("fi = %+v\n", fi)
	fmt.Printf("fi.IsDir() = %+v\n", fi.IsDir())
}
