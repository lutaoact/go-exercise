package main

import (
	"encoding/json"
	"fmt"
)

type foo struct {
	Message    string
	Ports      []int
	ServerName string
}

func newFoo() (*foo, error) {
	return &foo{
		Message:    "foo loves bar",
		Ports:      []int{80},
		ServerName: "Foo",
	}, nil
}

func main() {
	res, err := newFoo()

	out, err := json.Marshal(res)
	if err != nil {
		panic("xxxx")
	}

	fmt.Printf("string(out) = %+v\n", string(out))
}

/*
append(, value)
break
chan
case value:
const NAME Type = 0
const (
	NAME Type = iota

)
continue
defer func()

defer func()
defer func() {

}()

defer func() {
	if err := recover(); err != nil {

	}
}()

* This program is free software; you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation; either version 2 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program; if not, see <http://www.gnu.org/licenses/>.
*
* Copyright (C) Author,
import (
	"package"
)

type Interface interface {
}
if condition {

}
else {

}
if err := ; err != nil {

}
if err != nil {
	log.Fatal(err)
}
if err := xxxx; err != nil {

}

if err != nil {
	log.Fatal(err)
}
if err != nil {
	return err
}
if err != nil {
	return nil, err
}
if err != nil {
	panic()
}

if err != nil {
	t.Fatal(err)
}
if err != nil {

	return
}
`json:"json"`
`yaml:"yaml"`

fallthrough
for  {

}
for xxx := 0; xxx < N; xxx++ {
	xxx =
}
for k, v := range  {

}
func main() {
}


*/
