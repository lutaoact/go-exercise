package main

import (
	"fmt"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Pages            int    `json:"pages" bson:"pages"`
}

func NewBook(name string, pages int) *Book {
	return &Book{
		Name:  name,
		Pages: pages,
	}
}

func init() {
	// Setup mgm default config
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	fmt.Println(err)
}

func main() {
	book := NewBook("Pride and Prejudice", 345)

	// Make sure pass the model by reference.
	err := mgm.Coll(book).Create(book)
	fmt.Println(err)
}
