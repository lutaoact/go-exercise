package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func generator(done <-chan interface{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func multiply(done <-chan interface{}, in <-chan int, multiplier int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * multiplier:
			case <-done:
				return
			}
		}
	}()
	return out
}

func add(done <-chan interface{}, in <-chan int, additive int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n + additive:
			case <-done:
				return
			}
		}
	}()
	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// out的长度与nums相同，所以可以容纳所有的元素，但如果这样的话，其实就失去pipeline的内涵了，不过是从切片转成了channel，所以第一种才是正确的姿势
func gen2(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func sq5(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func merge2(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		for n := range c {
			select {
			case out <- n:
			case <-done:
			}
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func merge5(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func sqMain() {
	c := gen(2, 3)
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func sqMain2() {
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n)
	}
}

func sqMain3() {
	in := gen(2, 3, 4, 5, 6)
	c1 := sq(in)
	c2 := sq(in)

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}

func sqMain4() {
	in := gen(2, 3)
	c1 := sq(in)
	c2 := sq(in)

	done := make(chan struct{}, 2)
	out := merge2(done, c1, c2)
	fmt.Println(<-out)

	done <- struct{}{}
	done <- struct{}{}
}

func sqMain5() {
	done := make(chan struct{})
	defer close(done)

	in := gen(2, 3, 4, 5, 6)
	c1 := sq5(done, in)
	c2 := sq5(done, in)
	c3 := sq5(done, in)

	out := merge5(done, c1, c2, c3)
	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	return m, err
}

func md5DigestMain() {
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

func main() {
	//sqMain()
	//sqMain2()
	//sqMain3()
	//sqMain4()
	//sqMain5()
	//md5DigestMain()
	ppMain1()
}

func ppMain1() {
	done := make(chan interface{})
	nums := []int{1, 2, 3, 4}
	out := add(done, multiply(done, generator(done, nums...), 5), 3)
	for n := range out {
		fmt.Println(n)
	}
	close(done)
}
