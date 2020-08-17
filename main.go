package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	g := gin.Default()
	g.POST("/", RunCmd)
	g.GET("/healthCheck", func(ctx *gin.Context) {
		ctx.JSON(200,gin.H{
			"status":"ok",
		})
	})
	_ = g.Run(":9000")
	log.Println("alertmanager webhook start success!!!!")
}

