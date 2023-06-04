package routes

import (
	"context"
	"klikdaily-databoard/handlers"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteAPI(g *gin.Engine, parentCtx context.Context, db *gorm.DB) {
	repo := usecases.InitRepository(db)
	admin := handlers.InitVersionOneAdminHandler(repo)

	g.POST("/admins", admin.CreateAdmin)
	g.GET("/admins", admin.GetAdmins)
	g.GET("/admins/:id", admin.GetAdminById)
}
