package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Name    string
	Ports   []int
	Enabled bool
}

func main() {
	//timeTick()
	//timeParse()
	//timeFormat()
	fmt.Println(GetEndOfDay(time.Now()))
}

func GetEndOfDay(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("util.GetEndOfDay load location: %+v", err)
	}
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 23, 59, 59, 0, loc)
}

func timeFormat() {
	t := time.Now()
	fmt.Println(t.Format("2006.01.02 15:04"))
	fmt.Println(t.UTC().Format("2006.01.02 15:04"))
}

func timeTick() {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	fmt.Println(time.Now())
	for {
		select {
		case <-t.C:
			fmt.Println(time.Now())
		}
	}
}

func timeParse() {
	a, err := time.Parse("2006-01-02T15:04:05Z07:00", "2019-02-01T15:35:12.326+08:00")
	fmt.Println(a, err)
}

func TestTimeEqual(t *testing.T) {
	t.Log(time.RFC3339)
	assert := assert.New(t)
	t1, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(err)
	t2, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(err)
	assert.Equal(t1, t2, "t1 == t2")
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
