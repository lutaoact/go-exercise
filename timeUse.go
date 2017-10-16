package main

import (
	"fmt"
	"log"
	"time"
)

type Foo struct {
	Name    string
	Ports   []int
	Enabled bool
}

func main() {
	fmt.Println(time.RFC3339)
	t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("t = %+v\n", t)
}

/*
func typeOfTime() {
	t := time.Now()
	fmt.Println(reflect.TypeOf(t))

	latency := time.Since(t)
	fmt.Println(reflect.TypeOf(latency))
	fmt.Println(latency)

	foo := Foo{Name: "gopher", Ports: []int{80, 443}, Enabled: true}
	fmt.Printf("foo = %+v\n", foo)
}

func parseTime() {
	value := "Thu, 05/19/11"
	//	value := "Thu, 05/19/11, 10:47PM"
	// Writing down the way the standard time would look like formatted our way
	layout := "Mon, 01/02/06, 03:04PM"
	t, _ := time.Parse(layout, value)
	fmt.Println(t)
}

func formatTime() {
	//	t := time.SecondsToLocalTime(1305861602)
	//	var t time.Time
	//	t.ZoneOffset = -4 * 60 * 60
	//	fmt.Println(t.Format("2006-01-02 15:04:05 -0700"))
}
*/
