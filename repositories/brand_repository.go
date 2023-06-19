package repositories

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/middleware"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BrandRepositoryInterface interface {
	CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[any]
	GetBrandById(id string) chan RepositoryResult[any]
	GetBrands(br models.BrandRequest) chan RepositoryResult[any]
	GetPreviousBrand() string
	UpdateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand]
}

type brandRepository struct {
	db *gorm.DB
}

func InitBrandRepository(db *gorm.DB) BrandRepositoryInterface {
	return &brandRepository{
		db,
	}
}

func (b *brandRepository) GetBrands(br models.BrandRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		brands := []models.Brand{}
		brandCount := int64(0)
		page := br.Page
		limit := br.Limit
		b.db.Count(&brandCount)
		offset := (page - 1) * limit
		if err := b.db.Offset(int(offset)).Limit(int(limit)).Find(&brands).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       brands,
			Error:      nil,
			Message:    "Success Get List Brand",
			StatusCode: http.StatusOK,
		}
	}()
	return result
}

func (b *brandRepository) GetBrandById(id string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		brand := models.Brand{}
		err := b.db.First(&brand, id).Error
		if err == gorm.ErrRecordNotFound {
			result <- RepositoryResult[any]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Brand not found",
			}
			return
		} else if err != nil {
			result <- RepositoryResult[any]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusInternalServerError,
				Message:    "Error: " + err.Error(),
			}
			return
		} else if brand.Name == "" {
			result <- RepositoryResult[any]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Brand not found",
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       brand,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success get brand with id: " + brand.ID,
		}
	}()
	return result
}

func (b *brandRepository) GetPreviousBrand() string {
	latestId := ""
	if err := b.db.Model(&models.Brand{}).Select("id").Order("id desc").Limit(1).Scan(&latestId).Error; err != nil {
		return "error " + err.Error()
	}
	return latestId
}

func (b *brandRepository) CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	userName := middleware.ExtractNameFromToken(tokenString)
	// get by id first ? for create new id based on previous id
	latestId := b.GetPreviousBrand()
	if latestId == "" {
		latestId = "0"
	}
	latestIdInt, _ := strconv.Atoi(latestId)
	nextId := helper.GenerateNextID(latestIdInt)
	go func() {
		// validate brandRequest is created_by / updated_by superadmin role or adminGudang role or warehouseLeader role ?
		brand := models.Brand{
			Name:      br.Name,
			ID:        nextId,
			Status:    "active",
			CreatedBy: userName,
			UpdatedBy: userName,
		}
		b.db.Create(&brand)
		result <- RepositoryResult[any]{
			Data:       brand,
			Error:      nil,
			StatusCode: http.StatusCreated,
			Message:    "Success Create Brand",
		}
	}()
	return result
}

func (b *brandRepository) UpdateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand] {
	result := make(chan RepositoryResult[models.Brand])
	userName := middleware.ExtractNameFromToken(tokenString)
	go func() {
		brand := models.Brand{
			Name:      br.Name,
			ID:        br.ID,
			Status:    "active",
			CreatedBy: br.CreatedBy,
			UpdatedBy: userName,
		}
		b.db.Save(&brand)
		result <- RepositoryResult[models.Brand]{
			Data:       brand,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Update Brand Success",
		}
	}()
	return result
}
