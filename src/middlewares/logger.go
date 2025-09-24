package middlewares

import "github.com/gin-gonic/gin"

// Register attaches standard middleware used across the service.
func Register(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}
