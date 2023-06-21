package repositories

import (
	"fmt"
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any]
	GetPreviousId() string
}

type productRepository struct {
	db *gorm.DB
}

func InitProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &productRepository{
		db,
	}
}

func (r *productRepository) GetPreviousId() string {
	latestId := ""
	if err := r.db.Model(&models.Product{}).Select("id").Order("id desc").Limit(1).Scan(&latestId).Error; err != nil {
		return "error " + err.Error()
	}
	return latestId
}

func (r *productRepository) CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	latestID := r.GetPreviousId()
	latestOnlyId := helper.SplitProductID(latestID)
	latestIdInt, err := strconv.Atoi(latestOnlyId)
	if err != nil {
		result <- RepositoryResult[any]{
			Data:       nil,
			Error:      err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	productId, latestIdString := helper.GenerateProductID(latestIdInt)
	userName := helper.ExtractUserIDFromToken(tokenString)
	fmt.Println("latestID: " + latestID + "latest id string: " + latestIdString)
	fmt.Println(latestIdInt)
	go func() {
		product := models.Product{
			ID:         productId,
			Name:       pr.Name,
			BrandId:    pr.BrandId,
			Status:     "active",
			CategoryId: pr.CategoryId,
			CreatedBy:  userName,
			UpdatedBy:  userName,
			SKU:        pr.BrandId + pr.CategoryId + latestIdString,
		}
		if err := r.db.Create(&product).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    "Error: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Create Product Success",
			StatusCode: http.StatusCreated,
		}
	}()
	return result
}
