package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	pipe()
}

func pipe() {
	r, w := io.Pipe()

	go func() {
		fmt.Fprint(w, "some text to be read\n")
		w.Close()
	}()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	fmt.Print(buf.String())
}
