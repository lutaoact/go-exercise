package main

import (
	"math"
	"sync"
	"time"
)

func doAllocate(nKB int, wg *sync.WaitGroup) {
	var slice []byte
	for i := 0; i < nKB; i++ {
		t := make([]byte, 1024) // 1KB
		slice = append(slice, t...)
		// 大约会执行 50 秒，方便观察内存增长
		time.Sleep(time.Millisecond)
	}
	wg.Done()
	println("doAllocate done")
}
func doIdleAdd(n int64, wg *sync.WaitGroup) {
	var res int64
	for i := int64(0); i < n; i++ {
		res += i
	}
	wg.Done()
	println("doIdleAdd done")
}
func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) // needed before go 1.5
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go doAllocate(50*1024, wg) // 程序运行时最多分配 50MB-100MB 内存, 防止影响宿主机
	t := int64(math.Pow(10, 11))
	go doIdleAdd(t, wg) // 执行加法，大约会执行 30 秒，方便观察运行情况
	wg.Wait()
}
