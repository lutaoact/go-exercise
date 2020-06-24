package main

import "fmt"

func main() {
	//lenOfNilSlice()
	//sliceSelf()
	//slice()
	appendNilSlice()
}

func appendNilSlice() {
	var s []int = nil
	fmt.Println(s)
	s = append(s, 1)
	fmt.Println(s)

	var s2 = make(map[string][]string)
	fmt.Println(s2["hello"])
	a := s2["hello"]
	s2["hello"] = append(s2["hello"], "1")
	//s2["hello"] = append(a, "1")
	//s2["hello"] = make([]string, 0)
	fmt.Println(a)
	fmt.Println(s2)
}

func slice() {
	s := make([]int, 4, 8)
	s[0] = 5
	s[1] = 6
	s[2] = 7
	s[3] = 8

	s2 := s[1:3]
	fmt.Println(s2)
	fmt.Printf("cap(s2) = %+v\n", cap(s2))
	s2 = s2[:cap(s2)]
	fmt.Println(s2)
	fmt.Printf("cap(s2) = %+v\n", cap(s2))
}

func lenOfNilSlice() {
	var s []byte
	fmt.Println(len(s))
}

func sliceSelf() {
	s := &[]int{1, 2, 3, 4}
	*s = (*s)[0:3]
	fmt.Println(s)
}
