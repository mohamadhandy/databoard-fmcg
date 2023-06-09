package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
)

type CategoryHandlerInterface interface {
	CreateCategory(c *gin.Context)
	GetCategories(c *gin.Context)
}

type categoryHandler struct {
	categoryUseCase usecases.CategoryUseCaseInterface
}

func InitCategoryHandler(u usecases.CategoryUseCaseInterface) CategoryHandlerInterface {
	return &categoryHandler{
		categoryUseCase: u,
	}
}

func (ch *categoryHandler) CreateCategory(c *gin.Context) {
	categoryRequest := models.CategoryRequest{}
	authHeader := c.GetHeader("Authorization")
	c.BindJSON(&categoryRequest)
	categoryResult := ch.categoryUseCase.CreateCategory(authHeader, categoryRequest)
	c.JSON(categoryResult.StatusCode, categoryResult)
}

func (ch *categoryHandler) GetCategories(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	categoryResult := ch.categoryUseCase.GetCategories(authHeader)
	c.JSON(categoryResult.StatusCode, categoryResult)
}
