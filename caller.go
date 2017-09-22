package main

import (
	"fmt"
	"runtime"
)

func main() {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fmt.Printf("skip = %v, pc = %v, file = %v, line = %v\n", skip, pc, file, line)
	}
	/* MAC 上的输出 */
	/*
		skip = 0, pc = 17387715, file = /Users/lutao/go/src/github.com/lutaoact/go-exercise/caller.go, line = 10
		skip = 1, pc = 16943532, file = /usr/local/Cellar/go/1.9/libexec/src/runtime/proc.go, line = 185
		skip = 2, pc = 17106768, file = /usr/local/Cellar/go/1.9/libexec/src/runtime/asm_amd64.s, line = 2337
	*/
}
