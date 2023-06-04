package routes

import (
	"context"
	"klikdaily-databoard/handlers"
	"klikdaily-databoard/middleware"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteAPI(g *gin.Engine, parentCtx context.Context, db *gorm.DB) {
	repo := usecases.InitRepository(db)
	admin := handlers.InitVersionOneAdminHandler(repo)
	auth := handlers.InitVersionOneAuthenticationHandler(repo)

	g.POST("/login", auth.Login)

	g.POST("/admins", middleware.AuthMiddleware(), admin.CreateAdmin)
	g.GET("/admins", middleware.AuthMiddleware(), admin.GetAdmins)
	g.GET("/admins/:id", middleware.AuthMiddleware(), admin.GetAdminById)
}
