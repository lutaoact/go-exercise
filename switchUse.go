package main

import (
	"fmt"
	"strings"
)

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

func main1() {
	fmt.Println(switchType(&AppErr{}))
	fmt.Println(switchType(&AppErr2{}))
	fmt.Println(switchType(&AppErr3{}))
}

func main2() {
	name := "report"
	switch {
	case "report", "trial_class", "sku":
		// TODO
		fmt.Println("proxy to telisruby")
	case strings.HasPrefix(name, "sku_v2_"):
		fmt.Println("sku_v2_")
	default:
		fmt.Println("default")
	}
}

func main() {
	main2()
}
