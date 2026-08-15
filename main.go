package main

import (
	"app/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	route := gin.Default()

	route.Run(":8080")
}
