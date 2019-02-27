package main

import "fmt"

type AppErr struct {
	Status  int         `json:"-"`
	Msg     string      `json:"msg"`
	Payload interface{} `json:"payload"`
}

type AppErr2 struct{}

func (e *AppErr2) Error() string {
	return "xxx"
}

type AppErr3 struct{}

func switchType(e interface{}) int {
	switch body := e.(type) {
	case *AppErr:
		fmt.Printf("body = %+v\n", body)
		return 1
	case error:
		fmt.Printf("body = %+v\n", body)
		return 2
	default:
		return 3
	}
}

func main() {
	fmt.Println(switchType(&AppErr{}))
	fmt.Println(switchType(&AppErr2{}))
	fmt.Println(switchType(&AppErr3{}))
}
