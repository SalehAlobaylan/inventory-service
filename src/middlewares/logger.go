package middlewares

import "github.com/gin-gonic/gin"
// To do: make it real logging and enhance it
// Register attaches standard middleware used across the service.
func Register(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}
