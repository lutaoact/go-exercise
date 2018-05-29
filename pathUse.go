package main

import (
	"fmt"
	"path"
	"strings"
)

func main() {
	subPath := "var/log/xxxx"
	p := path.Join("tmp", subPath)
	fmt.Printf("p = %+v\n", p)

	subPath = "/var/log/xxxx/"
	p = path.Join("tmp", subPath)
	fmt.Printf("p = %+v\n", p)

	subPath = "/"
	p = path.Join("", subPath)
	fmt.Printf("p = %+v\n", p)

	subPath = "/hello/"
	fmt.Printf("p = %+v\n", ossPath(subPath))
}

func ossPath(path string) string {
	return strings.TrimLeft(strings.TrimRight("/xxxx", "/")+path, "/")
}
