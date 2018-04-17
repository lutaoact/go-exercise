package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

func main() {
	cpuProfile := "/data/backup/profile"
	startCPUProfile(&cpuProfile)

	for i := 0; i < 1000000; i++ {
		fmt.Println("hello")
	}
	stopCPUProfile(&cpuProfile)
}

func startCPUProfile(cpuProfile *string) {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create cpu profile output file: %s",
				err)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Can not start cpu profile: %s", err)
			f.Close()
			return
		}
	}
}

func stopCPUProfile(cpuProfile *string) {
	if *cpuProfile != "" {
		pprof.StopCPUProfile() // 把记录的概要信息写到已指定的文件
	}
}
