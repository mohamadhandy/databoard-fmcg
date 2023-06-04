package main

import (
	"context"
	"klikdaily-databoard/config"
	"klikdaily-databoard/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// router.POST("/users", authMiddleware, userController.RegisterUser)
	// router.POST("/login", userController.LoginUser)
	// router.GET("/users", authMiddleware, userController.GetAllUsers)
	// router.GET("/users/:user_id", authMiddleware, userController.GetUsersByID)
	// router.PUT("/users/:user_id", authMiddleware, userController.UpdateUser)
	// router.DELETE("/users/:user_id", authMiddleware, userController.DeleteUser)

	routes.RouteAPI(router, context.Background(), config.NewConnection())
	router.Run("localhost:9000")
}
