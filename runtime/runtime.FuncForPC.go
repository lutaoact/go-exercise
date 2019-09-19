package main

import (
	"fmt"
	"runtime"
)

func main() {
	for skip := 0; ; skip++ {
		pc, _, _, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		p := runtime.FuncForPC(pc)
		file, line := p.FileLine(0)

		fmt.Printf("skip = %v, pc = %v\n", skip, pc)
		fmt.Printf("  file = %v, line = %d\n", file, line)
		fmt.Printf("  entry = %v\n", p.Entry())
		fmt.Printf("  name = %v\n", p.Name())
	}
	// Output:
	// skip = 0, pc = 4198456
	//   file = caller.go, line = 8
	//   entry = 4198400
	//   name = main.main
	// skip = 1, pc = 4282882
	//   file = $(GOROOT)/src/pkg/runtime/proc.c, line = 179
	//   entry = 4282576
	//   name = runtime.main
	// skip = 2, pc = 4292528
	//   file = $(GOROOT)/src/pkg/runtime/proc.c, line = 1394
	//   entry = 4292528
	//   name = runtime.goexit
	pc := make([]uintptr, 1024)
	for skip := 0; ; skip++ {
		n := runtime.Callers(skip, pc)
		if n <= 0 {
			break
		}
		fmt.Printf("skip = %v, pc = %v\n", skip, pc[:n])
		for j := 0; j < n; j++ {
			p := runtime.FuncForPC(pc[j])
			file, line := p.FileLine(0)

			fmt.Printf("  skip = %v, pc = %v\n", skip, pc[j])
			fmt.Printf("    file = %v, line = %d\n", file, line)
			fmt.Printf("    entry = %v\n", p.Entry())
			fmt.Printf("    name = %v\n", p.Name())
		}
		break
	}
	// Output:
	// skip = 0, pc = [4307254 4198586 4282882 4292528]
	//   skip = 0, pc = 4307254
	//     file = $(GOROOT)/src/pkg/runtime/runtime.c, line = 315
	//     entry = 4307168
	//     name = runtime.Callers
	//   skip = 0, pc = 4198586
	//     file = caller.go, line = 8
	//     entry = 4198400
	//     name = main.main
	//   skip = 0, pc = 4282882
	//     file = $(GOROOT)/src/pkg/runtime/proc.c, line = 179
	//     entry = 4282576
	//     name = runtime.main
	//   skip = 0, pc = 4292528
	//     file = $(GOROOT)/src/pkg/runtime/proc.c, line = 1394
	//     entry = 4292528
	//     name = runtime.goexit
}
