package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qbox/stark/hms/api"
)

type Person struct {
	Name     string    `form:"name" binding:"required"`
	Address  string    `form:"address" binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		run2(c)
	})
	r.POST("/", func(c *gin.Context) {
		run3(c)
	})
	r.Run()
}

/*
curl -H "Content-type: application/json" -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16T15:04:05Z"}' "http://localhost:8080/"
curl -H "Content-type: application/json" -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16"}' "http://localhost:8080/"
curl -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16"}' "http://localhost:8080/"
curl -d 'name=appleboy&address=xyz&birthday=1992-03-16' http://localhost:8080/
*/
func run3(c *gin.Context) {
	var person Person
	//无法正确解析mimetype为application/json的post请求，time.Time格式的tag不处理
	//如果希望正确处理json，我觉得
	if err := c.ShouldBindWith(&person, binding.FormPost); err != nil {
		fmt.Printf("err = %+v\n", err)
		c.JSON(401, gin.H{"message": "pong"})
		return
	}
	fmt.Printf("person = %+v\n", person)
	log.Println(person.Name)
	log.Println(person.Address)
	log.Println(person.Birthday)
	c.String(200, "Success")
}

func run1(c *gin.Context) { //curl 'http://localhost:8080/ping?scope=1&scope=2'
	scopes := c.QueryArray("scope")
	fmt.Printf("scopes = %+v\n", scopes)
	c.JSON(200, gin.H{"message": "pong"})
}

func run2(c *gin.Context) { //curl 'http://localhost:8080/ping'
	err := api.NewHttpErr(401, "unauthorized", nil)
	fmt.Printf("err = %+v\n", err)
	api.E(c, err)
}
