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

type Rule struct {
	RuleID bson.ObjectId `json:"ruleID" bson:"_id"`
	Name   string        `json:"name"   bson:"name"`
}

func main() {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/hms")
	//	session, err := mgo.Dial("mongodb://127.0.0.1:28001,127.0.0.1:28002,127.0.0.1:28003/hms?replicaSet=kirk_rs1_dev")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("hello").C("people")
	c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	insertTest(c)

	var results []Person
	err = c.Find(bson.M{"name": "Ale"}).All(&results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone:", results[0])
}

func insertTest(c *mgo.Collection) {
	err := c.Insert(&Person{"Ale", "+1111111"}, &Person{"Cla", "+2222222"})
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		log.Fatal(err)
	}

	err = c.Insert(&Person{"Bbb", "+3333333"}, &Person{"Dxx", "+4444444"})
	if err != nil {
		log.Fatal(err)
	}
}
