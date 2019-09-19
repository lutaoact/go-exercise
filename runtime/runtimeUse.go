package main

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(path.Dir(filename))
	fmt.Println(strings.HasPrefix(runtime.Version(), "go1.10"))
}
