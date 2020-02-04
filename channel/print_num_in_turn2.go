package main

import "sync"

func main() {
	var done1 = make(chan struct{})
	var done2 = make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(100)

	go func() {
		for i := 1; i <= 100; i += 2 {
			<-done1
			println(i)
			wg.Done()
			done2 <- struct{}{}
		}
	}()

	go func() {
		for i := 2; i <= 100; i += 2 {
			<-done2
			println(i)
			wg.Done()
			done1 <- struct{}{} // 这里的最后一个值无法成功进入channel
		}
	}()

	done1 <- struct{}{}

	wg.Wait()

	close(done1)
	close(done2)
}
