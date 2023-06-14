package routes

import (
	"context"
	"klikdaily-databoard/handlers"
	"klikdaily-databoard/middleware"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RouteAPI(g *gin.Engine, parentCtx context.Context, db *gorm.DB, rdb *redis.Client) {
	repo := usecases.InitRepository(db, rdb)
	admin := handlers.InitVersionOneAdminHandler(repo)
	auth := handlers.InitVersionOneAuthenticationHandler(repo)
	brand := handlers.InitVersionOneBrandHandler(repo)
	category := handlers.InitVersionOneCategoryHandler(repo)

	g.POST("/login", auth.Login)

	g.POST("/admins", middleware.AuthMiddleware(), admin.CreateAdmin)
	g.GET("/admins", middleware.AuthMiddleware(), admin.GetAdmins)
	g.GET("/admins/:id", middleware.AuthMiddleware(), admin.GetAdminById)

	// brands
	g.POST("/brands", middleware.AuthMiddleware(), brand.CreateBrand)
	g.GET("/brands/:id", middleware.AuthMiddleware(), brand.GetBrandById)
	g.PUT("/brands/:id", middleware.AuthMiddleware(), brand.UpdateBrand)

	// categories
	g.POST("/categories", middleware.AuthMiddleware(), category.CreateCategory)
	g.GET("/categories", middleware.AuthMiddleware(), category.GetCategories)
}
