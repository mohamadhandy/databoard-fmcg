package handlers

import (
	"fmt"
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
)

type BrandHandlerInterface interface {
	CreateBrand(c *gin.Context)
}

type brandHandler struct {
	BrandUseCase usecases.BrandUsecaseInterface
}

func InitBrandHandler(u usecases.BrandUsecaseInterface) BrandHandlerInterface {
	return &brandHandler{
		BrandUseCase: u,
	}
}

func (b *brandHandler) CreateBrand(c *gin.Context) {
	// Get the token string from the Authorization header
	authHeader := c.GetHeader("Authorization")
	fmt.Println("auth", authHeader)
	br := models.BrandRequest{}
	c.BindJSON(&br)
	brandResult := b.BrandUseCase.CreateBrand(authHeader, br)
	c.JSON(brandResult.StatusCode, brandResult)
}
