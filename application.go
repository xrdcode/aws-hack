package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xrdcode/aws-hack/handler"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/upload", handler.Uploadimg)

	r.Run()
}
