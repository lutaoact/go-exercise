package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

// golang 中空值的定义
// The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
// 空的array、slice、map和字符串都会被忽略掉
type S1 struct {
	I1 int
	I2 int `json:",omitempty"`

	F1 float64
	F2 float64 `json:",omitempty"`

	S1 string
	S2 string `json:",omitempty"`

	B1 bool
	B2 bool `json:",omitempty"`

	Slice1 []int
	Slice2 []int `json:",omitempty"`
	Slice3 []int `json:",omitempty"`

	Map1 map[string]string
	Map2 map[string]string `json:",omitempty"`
	Map3 map[string]string `json:",omitempty"`

	O1 interface{}
	O2 interface{} `json:",omitempty"`
	O3 interface{} `json:",omitempty"`
	O4 interface{} `json:",omitempty"`
	O5 interface{} `json:",omitempty"`
	O6 interface{} `json:",omitempty"`
	O7 interface{} `json:",omitempty"`
	O8 interface{} `json:",omitempty"`

	P1 *int
	P2 *int               `json:",omitempty"`
	P3 *int               `json:",omitempty"`
	P4 *float64           `json:",omitempty"`
	P5 *string            `json:",omitempty"`
	P6 *bool              `json:",omitempty"`
	P7 *[]int             `json:",omitempty"`
	P8 *map[string]string `json:",omitempty"`
}

func main() {

	p3 := 0
	p4 := float64(0)
	p5 := ""
	p6 := false
	p7 := []int{}
	p8 := map[string]string{}

	s1 := S1{
		I1: 0,
		I2: 0,

		F1: 0,
		F2: 0,

		S1: "",
		S2: "",

		B1: false,
		B2: false,

		Slice1: []int{},
		Slice2: nil,
		Slice3: []int{},

		Map1: map[string]string{},
		Map2: nil,
		Map3: map[string]string{},

		O1: nil,
		O2: nil,
		O3: int(0),
		O4: float64(0),
		O5: "",
		O6: false,
		O7: []int{},
		O8: map[string]string{},

		P1: nil,
		P2: nil,
		P3: &p3,
		P4: &p4,
		P5: &p5,
		P6: &p6,
		P7: &p7,
		P8: &p8,
	}

	b, err := json.Marshal(s1)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)

	//Output:
	//{
	//	"I1": 0,
	//	"F1": 0,
	//	"S1": "",
	//	"B1": false,
	//	"Slice1": [],
	//	"Map1": {},
	//	"O1": null,
	//	"O3": 0,
	//	"O4": 0,
	//	"O5": "",
	//	"O6": false,
	//	"O7": [],
	//	"O8": {},
	//	"P1": null,
	//	"P3": 0,
	//	"P4": 0,
	//	"P5": "",
	//	"P6": false,
	//	"P7": [],
	//	"P8": {}
	//}
}
