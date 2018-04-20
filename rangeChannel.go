package main

import (
	"fmt"
	"strconv"
	"time"
)

func makeCakeAndSend(cs chan string, count int) {
	for i := 1; i <= count; i++ {
		cakeName := "Strawberry Cake " + strconv.Itoa(i)
		fmt.Println("in:", i)
		cs <- cakeName //send a strawberry cake
	}
	close(cs)
}

func receiveCakeAndPack(cs chan string) {
	for s := range cs {
		fmt.Println("Packing received cake: ", s)
	}
}

func main() {
	cs := make(chan string, 5)
	go makeCakeAndSend(cs, 5)

	//这里sleep的时候，channel已经close了，但其中的数据尚未被消费
	time.Sleep(3 * 1e9)

	//这里才是去消费数据的地方
	go receiveCakeAndPack(cs)

	//sleep for a while so that the program doesn’t exit immediately
	for {
	}
}
