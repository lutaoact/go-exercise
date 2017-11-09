package main

import "fmt"

type user struct {
	name  string
	email string
}

/*
如果值作为接收者声明方法，调用时会使用这个值的一个副本来执行，如notify和changeEmail方法，调用并不会改变原始值。
如果指针作为接收者声明方法，调用之后，外面的值同样会被改变，如changeEmail2.
*/

func (u user) notify() {
	fmt.Println(u.name, u.email)
}

//值作为接收者，在函数内部改变属性，
func (u user) changeEmail(email string) {
	u.email = email
}

func (u *user) changeEmail2(email string) {
	u.email = email
}

func main() {
	bill := user{"Bill", "bill@email.com"}
	bill.notify()

	lisa := &user{"Lisa", "lisa@email.com"}
	lisa.notify()

	bill.changeEmail("bill@newdomail.com")
	bill.notify()

	lisa.changeEmail("lisa@newdomail.com")
	lisa.notify()

	bill.changeEmail2("bill@newdomail.com")
	bill.notify()

	lisa.changeEmail2("lisa@newdomail.com")
	lisa.notify()
}

/*
OUTPUT:
Bill bill@email.com
Lisa lisa@email.com
Bill bill@email.com
Lisa lisa@email.com
Bill bill@newdomail.com
Lisa lisa@newdomail.com
*/
