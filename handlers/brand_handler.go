package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BrandHandlerInterface interface {
	CreateBrand(c *gin.Context)
	GetBrandById(c *gin.Context)
	UpdateBrand(c *gin.Context)
	GetBrands(c *gin.Context)
}

type brandHandler struct {
	BrandUseCase usecases.BrandUsecaseInterface
}

func InitBrandHandler(u usecases.BrandUsecaseInterface) BrandHandlerInterface {
	return &brandHandler{
		BrandUseCase: u,
	}
}

func (b *brandHandler) GetBrands(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	brandRequest := models.BrandRequest{
		Page:  p,
		Limit: l,
	}
	brands := b.BrandUseCase.GetBrands(brandRequest)
	c.JSON(brands.StatusCode, brands)
}

func (b *brandHandler) CreateBrand(c *gin.Context) {
	// Get the token string from the Authorization header
	authHeader := c.GetHeader("Authorization")
	br := models.BrandRequest{}
	c.BindJSON(&br)
	brandResult := b.BrandUseCase.CreateBrand(authHeader, br)
	c.JSON(brandResult.StatusCode, brandResult)
}

func (b *brandHandler) GetBrandById(c *gin.Context) {
	inputId := c.Param("id")
	brandResult := b.BrandUseCase.GetBrandById(inputId)
	c.JSON(brandResult.StatusCode, brandResult)
}

func (b *brandHandler) UpdateBrand(c *gin.Context) {
	inputId := c.Param("id")
	authHeader := c.GetHeader("Authorization")
	br := models.BrandRequest{
		ID: inputId,
	}
	c.BindJSON(&br)
	brandResult := b.BrandUseCase.UpdateBrand(authHeader, br)
	c.JSON(brandResult.StatusCode, brandResult)
}
