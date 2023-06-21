package repositories

import (
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
	latestOnlyId := ""
	var latestIdInt int
	go func() {
		if latestID != "" {
			latestOnlyId = helper.SplitProductID(latestID)
			var errConvert error
			latestIdInt, errConvert = strconv.Atoi(latestOnlyId)
			if errConvert != nil {
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errConvert,
					Message:    errConvert.Error(),
					StatusCode: http.StatusInternalServerError,
				}
				// if err occurs from Atoi then return
				return
			}
		} else if latestID == "" {
			var errConvert error
			latestIdInt, errConvert = strconv.Atoi("0") // why "0"? because first time create product.
			if errConvert != nil {
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errConvert,
					Message:    errConvert.Error(),
					StatusCode: http.StatusInternalServerError,
				}
				return
			}
		}
		productId, latestIdString := helper.GenerateProductID(latestIdInt)
		userName := helper.ExtractUserIDFromToken(tokenString)
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
