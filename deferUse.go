package main

import "fmt"

func main() {
	a := deferTest()
	fmt.Printf("a = %+v\n", a)

	b := f1()
	fmt.Printf("b = %+v\n", b)

	c := f2()
	fmt.Printf("c = %+v\n", c)

	d := f3()
	fmt.Printf("d = %+v\n", d)
}

//如果在defer语句执行之前，函数就return了，那么defer的函数是没有机会执行的
//所以，deferTest函数的返回值就是2
func deferTest() int {
	a := 2
	return a

	defer func() {
		a *= 2
		fmt.Println(1)
	}()
	return 3
}

//我认为会返回0，实际返回1
func f1() (result int) {
	defer func() {
		result++
	}()
	return 0
}

//我认为会返回10，实际返回5
func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

//我认为会返回5，实际返回1
func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

//f1, f2, f3的调用结果，我刚开始全都想错了，下面来解释一下正确的理论：
//函数返回的过程是这样的：先给返回值赋值，然后调用defer表达式，最后才是返回到调用函数中。
/*
f1中，return的时候，result被赋值为0，然后调用defer，result自加得到1，所以返回值为1.
f2中，return的时候，r被赋值为t的值，也就是5，然后调用defer，t变为10，但并不影响r的值，所以返回值为5
f3中，return的时候，r被赋值为1，然后调用defer，这里要注意，r作为参数传递，于是函数中使用的是局部变量r，它的值改变并不会影响外面的值，所以返回1
*/
