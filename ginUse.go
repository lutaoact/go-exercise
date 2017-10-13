package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qbox/stark/hms/api"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		run2(c)
	})
	r.Run()
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
