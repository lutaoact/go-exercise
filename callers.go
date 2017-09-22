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

	pc := make([]uintptr, 1024)
	for skip := 0; ; skip++ {
		n := runtime.Callers(skip, pc) //该函数返回写入到 pc 切片中的项数(受切片的容量限制)
		if n <= 0 {
			break
		}
		fmt.Printf("skip = %v, pc = %v\n", skip, pc[:n])
	}
	// Output:
	// skip = 0, pc = [4304486 4198562 4280114 4289760]
	// skip = 1, pc = [4198562 4280114 4289760]
	// skip = 2, pc = [4280114 4289760]
	// skip = 3, pc = [4289760]
}

/*
输出新的 pc 长度和 skip 大小有逆相关性. skip = 0 为 runtime.Callers 自身的信息.

skip = 0, pc = 17387739, file = /Users/lutao/go/src/github.com/lutaoact/go-exercise/callers.go, line = 10
skip = 1, pc = 16943532, file = /usr/local/Cellar/go/1.9/libexec/src/runtime/proc.go, line = 185
skip = 2, pc = 17106768, file = /usr/local/Cellar/go/1.9/libexec/src/runtime/asm_amd64.s, line = 2337
skip = 0, pc = [16804753 17388146 16943533 17106769]
skip = 1, pc = [17388146 16943533 17106769]
skip = 2, pc = [16943533 17106769]
skip = 3, pc = [17106769]

https://github.com/Unknwon/gcblog/blob/master/content/04-go-caller.md
输出跟文章中的不一样，为什么会相差1呢？
*/
