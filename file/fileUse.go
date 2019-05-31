package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	mainLstat()
}

func mainLstat() {
	info, err := os.Lstat("/Users/taolu/exercise/")
	fmt.Printf("info = %+v\n", info)
	fmt.Println(err)
	fmt.Println(info.Mode())
	fmt.Println(info.Mode().IsRegular())
}

func mainGlob() {
	filepath.Glob()

}
