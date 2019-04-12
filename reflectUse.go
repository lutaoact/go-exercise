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

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

func main() {
	//testValueOf()
	//func1()
	convert()
}

func testValueOf() {
	SwaggerInfo = swaggerInfo{
		Version: "1.11",
		Host:    "xxxx",
	}
	v := reflect.ValueOf(SwaggerInfo)

	fmt.Printf("v = %+v\n", v)
}

func func1() {
	var person Person

	typ := reflect.TypeOf(&person)
	fmt.Printf("typ = %+v\n", typ)
	typelem := reflect.TypeOf(&person).Elem()
	fmt.Printf("typelem = %+v\n", typelem)
	val := reflect.ValueOf(&person).Elem()
	fmt.Printf("val = %+v\n", val)
	fmt.Printf("typ.NumField() = %+v\n", typelem.NumField())

	fmt.Println(reflect.TypeOf((*Person)(nil)).Elem())

	s := "hello"
	fmt.Println(reflect.TypeOf(&s).Elem())

	//	for i := 0; i < typ.NumField(); i++ {
	//		typeField := typ.Field(i)
	//		structField := val.Field(i)
	//		fmt.Printf("typeField = %+v\n", typeField)
	//		fmt.Printf("structField = %+v\n", structField)
	//		fmt.Printf("typeField.Tag.Get(time_format) = %+v\n", typeField.Tag.Get("time_format"))
	//	}
}

func convert() {
	var num float64 = 1.2345

	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)

	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer)
	fmt.Println(convertValue)
}
