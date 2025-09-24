package routes

import (
	"github.com/gin-gonic/gin"

	"inventory-service/src/controllers"
)

// RegisterRoutes wires all inventory endpoints to the provided Gin engine.
func RegisterRoutes(router *gin.Engine) {
	inventory := router.Group("/inventory")
	{
		inventory.GET("", controllers.GetItems)
		inventory.POST("", controllers.CreateItem)
		inventory.GET("/:id", controllers.GetItemByID)
		inventory.PUT("/:id", controllers.UpdateItem)
		inventory.DELETE("/:id", controllers.DeleteItem)
	}
}
