package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/api", ApiHandler)
	r.GET("/", StaticHandler)
	r.Run() // listen and serve on 0.0.0.0:8080
}

type User struct {
	Name string
	Age  int
}

func StaticHandler(c *gin.Context) {
	user := User{
		Name: "example",
		Age:  20,
	}
	c.HTML(http.StatusOK, "views/index.tmpl", gin.H{
		"Name": user.Name,
		"Age":  user.Age,
	})
}

func ApiHandler(c *gin.Context) {
	user := User{
		Name: "example",
		Age:  20,
	}
	c.JSON(200, gin.H{
		"Name": user.Name,
		"Age":  user.Age,
	})
}