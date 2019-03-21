package main

func main() {
	var done1 = make(chan struct{}, 1) //需要有一个channel是带buffer的，这个channel负责接最后一个值
	var done2 = make(chan struct{})

	var done = make(chan struct{})

	go func() {
		for i := 1; i <= 10; i += 2 {
			<-done1
			println(i)
			done2 <- struct{}{}
		}
	}()

	go func() {
		for i := 2; i <= 10; i += 2 {
			<-done2
			println(i)
			done1 <- struct{}{} // 这里接最后一个值，然后退出循环，如果没有buffer，就会阻塞
		}
		done <- struct{}{}
	}()

	done1 <- struct{}{}

	<-done
}
