package main

import (
	"demo/router"

	"github.com/gin-gonic/gin"
)

func main() {

	router1 := gin.Default()

	router.Userroutes(router1)

	router1.Run("localhost:8080")
}
