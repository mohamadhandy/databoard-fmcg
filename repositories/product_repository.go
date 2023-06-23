package repositories

import (
	"errors"
	"fmt"
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	CreateProduct(pr models.ProductRequest, tokenString string) chan RepositoryResult[any]
	GetProductById(id string) chan RepositoryResult[any]
	GetProducts(productReq models.ProductRequest) chan RepositoryResult[any]
	UpdateProduct(tokenString string, productReq models.ProductRequest) chan RepositoryResult[any]
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

func (r *productRepository) GetProducts(productReq models.ProductRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		products := []models.ProductResponse{}
		page := productReq.Page
		limit := productReq.Limit
		query := fmt.Sprintf(`SELECT p2.id, p2.name, c2.name as category_name,
		b2.name as brand_name, p2.status, p2.sku,
		p2.created_at, p2.updated_at,
		p2.created_by, p2.updated_by
		FROM "Product" as p2
		JOIN "Category" as c2 ON p2.category_id = c2.id
		JOIN "Brand" as b2 ON p2.brand_id = b2.id
		LIMIT %v OFFSET (%v - 1) * %v`, limit, page, limit)
		if err := r.db.Raw(query).Scan(&products).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       products,
			Error:      nil,
			Message:    "Success Get Products",
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (r *productRepository) GetPreviousId() string {
	latestID := ""
	if err := r.db.Model(&models.Product{}).Select("id").Order("id desc").Limit(1).Scan(&latestID).Error; err != nil {
		return "error " + err.Error()
	}
	return latestID
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

func (r *productRepository) UpdateProduct(tokenString string, productReq models.ProductRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		userName := helper.ExtractUserIDFromToken(tokenString)
		product := models.Product{
			Name:       productReq.Name,
			ID:         productReq.ID,
			BrandId:    productReq.BrandId,
			CategoryId: productReq.CategoryId,
			UpdatedBy:  userName,
			Status:     "active",
			CreatedBy:  productReq.CreatedBy,
			CreatedAt:  productReq.CreatedAt,
			SKU:        productReq.SKU,
		}
		tx := r.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      errors.New("panic occured"),
					StatusCode: http.StatusInternalServerError,
					Message:    "An unexpected token",
				}
				return
			}
		}()

		err := tx.Transaction(func(tx *gorm.DB) error {
			if err := r.db.Save(&product).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			tx.Rollback()
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       product,
			Error:      nil,
			Message:    "Success Update Product",
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
		if strings.Contains(latestID, "error") {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      errors.New(latestID),
				StatusCode: http.StatusInternalServerError,
				Message:    latestID,
			}
			return
		}
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
