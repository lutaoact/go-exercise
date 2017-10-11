package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		scopes := c.QueryArray("scope") //curl 'http://localhost:8080/ping?scope=1&scope=2'
		fmt.Printf("scopes = %+v\n", scopes)
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run()
}
