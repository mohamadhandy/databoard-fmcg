package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandlerInterface interface {
	CreateProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	GetProducts(c *gin.Context)
}

type productHandler struct {
	productUseCase usecases.ProductUseCaseInterface
}

func InitProductHandler(u usecases.ProductUseCaseInterface) ProductHandlerInterface {
	return &productHandler{
		productUseCase: u,
	}
}

func (h *productHandler) GetProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	productReq := models.ProductRequest{
		Page:  pageInt,
		Limit: limitInt,
	}
	result := h.productUseCase.GetProducts(productReq)
	c.JSON(result.StatusCode, result)
}

func (h *productHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	product := h.productUseCase.GetProductById(id)
	c.JSON(product.StatusCode, product)
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	productRequest := models.ProductRequest{}
	authHeader := c.GetHeader("Authorization")
	c.BindJSON(&productRequest)
	productResult := h.productUseCase.CreateProduct(authHeader, productRequest)
	c.JSON(productResult.StatusCode, productResult)
}
