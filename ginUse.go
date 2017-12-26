package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name     string    `form:"name"     json:"name"                                           binding:"required"`
	Address  string    `form:"address"  json:"address"                                        binding:"required"`
	Birthday time.Time `form:"birthday" json:"birthday" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

func main() {
	r := gin.Default()

	//curl -v 'http://127.0.0.1:8080/context'
	/*
		func (c *Context) Value(key interface{}) interface{} {
			if key == 0 {
				return c.Request
			}
			if keyAsString, ok := key.(string); ok {
				val, _ := c.Get(keyAsString)
				return val
			}
			return nil
		}
	*/
	r.GET("/context", func(c *gin.Context) {
		logrus.Info(c.Value(0)) //c.Request字段，类型为*http.Request
		logrus.Info(context.Background().Value(0))
		c.String(200, "context")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	//curl -v -H 'Content-Type: application/json' -d '{"name":"hello","address":"world"}' 'http://127.0.0.1:8080/hook'
	r.POST("/hook", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("err = %+v\n", err)
			c.String(500, "failure")
			return
		}
		fmt.Println(string(data))
		c.String(200, "Success")
	})

	//curl -v -H 'Content-Type: application/json' -d '{"name":"hello","address":"world"}' 'http://127.0.0.1:8080/testHttpPost'
	r.POST("/testHttpPost", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("err = %+v\n", err)
			return
		}

		var p Person
		fmt.Println(string(data))
		err = json.Unmarshal(data, &p)
		if err != nil {
			fmt.Printf("err = %+v\n", err)
			return
		}
		fmt.Printf("p = %+v\n", p)

		lines := make([]string, 0)
		for k, v := range c.Request.Header {
			lines = append(lines, fmt.Sprintf("%s: %s", k, strings.Join(v, "")))
		}
		fmt.Println(strings.Join(lines, "\n"))

		c.JSON(200, &Person{"hello", "world", time.Now()})
	})

	//curl -v 'http://127.0.0.1:8080/testHeader'
	r.GET("/testHeader", func(c *gin.Context) {
		lines := make([]string, 0)
		for k, v := range c.Request.Header {
			lines = append(lines, fmt.Sprintf("%s: %s", k, strings.Join(v, "")))
		}
		fmt.Println(strings.Join(lines, "\n"))
		c.String(200, "Success")
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
	r.PUT("/updateRepo", func(c *gin.Context) {
		run4(c)
	})
	r.Run()
}

/*
updatedata='{"namespace":"lutaoact","name":"mynginx","logoUrl":"http://ozdy4xcj5.bkt.cloudappl.com/nginx_large.png","summary":"this is a short summary","description":"nginx is great and great and great.","origin":"docker","labels":["database","os"],"tags":["1.0","1.1","latest"],"codeSource":"github","isPub":true}'
updatedata='{"namespace":"lutaoact","name":"mynginx","logoUrl":"http://ozdy4xcj5.bkt.cloudappl.com/nginx_large.png","summary":"this is a short summary","description":"nginx is great and great and great.","codeSource":"github","isPub":true}'

curl -v -X PUT "Content-Type: application/json" -d "$updatedata" 'http://127.0.0.1:8080/updateRepo'
*/
func run4(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "%s", err)
		return
	}
	var repo map[string]interface{}
	err = json.Unmarshal(data, &repo)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "%s", err)
		return
	}
	fmt.Printf("repo = %+v\n", repo)
	fmt.Printf("len(repo.tags) = %+v\n", len(repo["tags"].([]interface{})))
	//	fmt.Printf("repo.tags = %+v\n", repo["tags"].([]string)[0])
	c.JSON(http.StatusOK, repo["tags"].([]interface{}))
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
