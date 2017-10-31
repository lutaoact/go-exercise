package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name     string    `form:"name" binding:"required"`
	Address  string    `form:"address" binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
	})
	r.GET("/", func(c *gin.Context) {
		filter := GetBaseFilter(c)
		filter["xxx"] = "yyyy"
		fmt.Printf("filter = %+v\n", filter)
		c.String(200, "Success")
	})
	r.POST("/", func(c *gin.Context) {
		run3(c)
	})
	r.Run()
}

//curl 'http://localhost:8080?origin=gogogo&label=server&isCertified=true'
func GetBaseFilter(c *gin.Context) bson.M {
	filter := bson.M{}
	if c.Query("origin") != "" {
		filter["origin"] = c.Query("origin")
	}
	if c.Query("label") != "" {
		filter["labels"] = c.Query("label") //注意：lables是数据库中的数组字段
	}
	if c.Query("isCertified") == "true" {
		filter["isCertified"] = true
	}
	return filter
}

/*
curl -H "Content-type: application/json" -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16T15:04:05Z"}' "http://localhost:8080/"
curl -H "Content-type: application/json" -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16"}' "http://localhost:8080/"
curl -d '{"name":"appleboy","address":"xyz","birthday":"1992-03-16"}' "http://localhost:8080/"
curl -d 'name=appleboy&address=xyz&birthday=1992-03-16' http://localhost:8080/
*/
/*
 * 无法正确解析mimetype为application/json的post请求，time.Time格式的tag不处理
 * 原理是这样的：
 *   post请求中的body数据是利用json的decoder直接解析的，time.Time类型实现了
 *   Marshaler和Unmarshaler接口，可以看到format字符串是固定的RFC3339，并没有
 *   识别time_format这种tag的地方
 *
 */
func run3(c *gin.Context) {
	var person Person
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

/*
func run2(c *gin.Context) { //curl 'http://localhost:8080/ping'
	err := api.NewHttpErr(401, "unauthorized", nil)
	fmt.Printf("err = %+v\n", err)
	api.E(c, err)
}
*/
