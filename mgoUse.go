package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("hello").C("people")

	err = c.Insert(&Person{"Ale", "+1111111"}, &Person{"Cla", "+2222222"})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Insert(&Person{"Bbb", "+3333333"}, &Person{"Dxx", "+4444444"})
	if err != nil {
		log.Fatal(err)
	}

	var results []Person
	err = c.Find(bson.M{"name": "Ale"}).All(&results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone:", results)
}
