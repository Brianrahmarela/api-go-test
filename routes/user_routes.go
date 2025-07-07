package routes

import (
	"api-go-test/controllers"
	"api-go-test/middleware"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	protected := router.Group("/")
	nonProtected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", userController.GetUsers)
	}

	nonProtected.GET("/users/:id", userController.GetUserByID)
}
