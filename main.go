package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.POST("/", RunCmd)
	_ = g.Run(":9000")
}

