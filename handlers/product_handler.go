package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
)

type ProductHandlerInterface interface {
	CreateProduct(c *gin.Context)
}

type productHandler struct {
	productUseCase usecases.ProductUseCaseInterface
}

func InitProductHandler(u usecases.ProductUseCaseInterface) ProductHandlerInterface {
	return &productHandler{
		productUseCase: u,
	}
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	productRequest := models.ProductRequest{}
	authHeader := c.GetHeader("Authorization")
	c.BindJSON(&productRequest)
	productResult := h.productUseCase.CreateProduct(authHeader, productRequest)
	c.JSON(productResult.StatusCode, productResult)
}
