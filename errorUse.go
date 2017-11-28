package main

import (
	"errors"
	"fmt"
)

func main() {
	testFunc := func() (err error) {
		a, err := 1, errors.New("a error") // a是新创建变量，err是被赋值
		//		if err != nil {
		//			return // 正确返回err
		//		}
		fmt.Printf("a = %+v\n", a)
		fmt.Printf("err = %+v\n", err)

		if b, err := 2, errors.New("b error"); err != nil { // 此刻if语句中err被重新创建
			return err // if语句中的err覆盖外面的err，导致编译
			//  错误 (err is shadowed during return。)
		} else {
			fmt.Println(b)
		}
		return
	}

	err := testFunc()
	fmt.Printf("err = %+v\n", err)
}
