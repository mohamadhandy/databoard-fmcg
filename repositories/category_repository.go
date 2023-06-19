package repositories

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/middleware"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	CreateCategory(tokenString string, cr models.CategoryRequest) chan RepositoryResult[any]
	GetCategories(tokenString string) chan RepositoryResult[any]
	GetPreviousCategoryId() string
}

type categoryRepository struct {
	db *gorm.DB
}

func InitCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepository{
		db,
	}
}

func (c *categoryRepository) GetPreviousCategoryId() string {
	latestId := ""
	if err := c.db.Model(&models.Category{}).Select("id").Order("id desc").Limit(1).Scan(&latestId).Error; err != nil {
		return "error " + err.Error()
	}
	return latestId
}

func (c *categoryRepository) GetCategories(tokenString string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		categoryRes := []models.Category{}
		if err := c.db.Find(&categoryRes).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       []models.Category{},
				Error:      err,
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       categoryRes,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success get Categories",
		}
	}()
	return result
}

func (c *categoryRepository) CreateCategory(tokenString string, cr models.CategoryRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	userName := middleware.ExtractNameFromToken(tokenString)
	latestId := c.GetPreviousCategoryId()
	if latestId == "" {
		latestId = "0"
	}
	latestIdInt, _ := strconv.Atoi(latestId)
	nextId := helper.GenerateNextIDCategory(latestIdInt)
	go func() {
		category := models.Category{
			ID:        nextId,
			CreatedBy: userName,
			UpdatedBy: userName,
			Name:      cr.Name,
		}
		c.db.Create(&category)
		result <- RepositoryResult[any]{
			Data:       category,
			StatusCode: http.StatusCreated,
			Error:      nil,
			Message:    "Create Category Success",
		}
	}()
	return result
}
