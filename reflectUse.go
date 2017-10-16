package main

import (
	"fmt"
	"reflect"
	"time"
)

type Person struct {
	Name     string    `form:"name" binding:"required"`
	Address  string    `form:"address" binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

func main() {
	var person Person

	typ := reflect.TypeOf(&person).Elem()
	fmt.Printf("typ = %+v\n", typ)
	val := reflect.ValueOf(&person).Elem()
	fmt.Printf("val = %+v\n", val)
	fmt.Printf("typ.NumField() = %+v\n", typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		fmt.Printf("typeField = %+v\n", typeField)
		fmt.Printf("structField = %+v\n", structField)
		fmt.Printf("typeField.Tag.Get(time_format) = %+v\n", typeField.Tag.Get("time_format"))
	}
}
