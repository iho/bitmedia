package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/users", createUser)
	router.GET("/users", listUsers)

	router.GET("/games", listUsers)
	router.GET("/games/:id/tats", listUsers)

	router.Run()
}
