package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

/*
I use 4 goroutines to sort integers, I use priority queue to merge result.
*/

type Item struct {
	partIndex int // belong to which part
	priority  int // value of item, need to sort
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func sortInts(partIndex int, s []int) []int {
	fmt.Printf("part%d = %+v\n", partIndex, s)
	sort.Ints(s)
	return s
}

func main() {
	fmt.Println("Please input a integer list, use space or new line to split:")
	fmt.Println("For example: 98 23 1 55 6 22 71 6 8 25 31 2")
	fmt.Println("Use Ctrl-D to stop in OSX or linux, Ctrl-Z in windows.")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	list := make([]int, 0)
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		list = append(list, x)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if len(list) < 4 {
		panic("length of list should be greater than 4")
	}

	// use math.Ceil
	step := int(math.Ceil(float64(len(list)) / 4))

	acceptor := make(chan []int)
	for i := 0; i < 4; i++ {
		part := make([]int, step)
		if i == 3 {
			part = list[i*step : len(list)]
		} else {
			part = list[i*step : (i+1)*step]
		}
		go func(i int, part []int) {
			acceptor <- sortInts(i, part)
		}(i, part)
	}

	parts := make(map[int][]int)
	for i := 0; i < 4; i++ {
		parts[i] = <-acceptor
	}
	//fmt.Println(parts)

	pq := make(PriorityQueue, 4)
	for i := 0; i < 4; i++ {
		pq[i] = &Item{
			partIndex: i,
			priority:  parts[i][0],
		}
		parts[i] = parts[i][1:]
	}
	heap.Init(&pq)

	results := make([]int, len(list))
	for i := 0; i < len(list); i++ {
		//fmt.Println(parts)
		item := heap.Pop(&pq).(*Item)
		results[i] = item.priority
		if len(parts[item.partIndex]) == 0 {
			continue
		}

		heap.Push(&pq, &Item{
			partIndex: item.partIndex,
			priority:  parts[item.partIndex][0],
		})
		parts[item.partIndex] = parts[item.partIndex][1:]
	}
	fmt.Println("sorted results:", results)
}
