package repositories

import (
	"errors"
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any]
	GetProductById(id string) chan RepositoryResult[any]
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

func (r *productRepository) GetProductById(id string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		product := models.Product{}
		if err := r.db.Where("id = ?", id).Find(&product).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		if product.Name == "" {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      errors.New("product not found"),
				Message:    "Product Not found",
				StatusCode: http.StatusNotFound,
			}
		}
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Success get product by id: " + id,
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (r *productRepository) CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		tx := r.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errors.New("panic occurred"),
					Message:    "An unexpected error occurred",
					StatusCode: http.StatusInternalServerError,
				}
			}
		}()

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
			return
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

		err = tx.Transaction(func(tx *gorm.DB) error {
			if err := r.db.Create(&product).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			tx.Rollback()
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
