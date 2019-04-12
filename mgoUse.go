package main

import (
	"fmt"
	"log"
	"time"

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

type Namespace struct {
	Name      string    `json:"name"      bson:"name"  index:"unique"`
	IsPub     bool      `json:"isPub"     bson:"isPub"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func main() {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/hms")
	//	session, err := mgo.Dial("mongodb://127.0.0.1:28001,127.0.0.1:28002,127.0.0.1:28003/hms?replicaSet=kirk_rs1_dev")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	TestNotFound(session)

	//TestAggregate(session)
}

func TestAggregate(session *mgo.Session) {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"namespace": "namespaceForTest",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$namespace",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	result := []bson.M{}
	session.DB("hms").C("images").Pipe(pipeline).All(&result)
	fmt.Println(result)
}

func TestDistinct(session *mgo.Session) {
	c := session.DB("hello").C("people")
	names := make([]Person, 0)
	err := c.Find(nil).Sort("phone").All(&names)
	fmt.Println(names, err)
}

func TestEnsureIndex(session *mgo.Session) {
	c := session.DB("hello").C("people")
	c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	insertTest(c)

	var results []Person
	err := c.Find(bson.M{"name": "Ale"}).All(&results)
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

func TestNotFound(session *mgo.Session) {
	c := session.DB("hms").C("namespaces")
	var ns Namespace
	err := c.Find(bson.M{"name": "lutaoact2"}).One(&ns)
	fmt.Printf("err = %+v\n", err)
	/*
		if err.Error() == "not found" {
			fmt.Println("not found 1")
		}
	*/
	if err == mgo.ErrNotFound {
		fmt.Println("not found 2")
	}
}
