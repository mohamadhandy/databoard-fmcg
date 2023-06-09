package repositories

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BrandRepositoryInterface interface {
	CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand]
	GetBrandById(id string) chan RepositoryResult[models.Brand]
	GetPreviousBrand() string
}

type brandRepository struct {
	db *gorm.DB
}

func InitBrandRepository(db *gorm.DB) BrandRepositoryInterface {
	return &brandRepository{
		db,
	}
}

func (b *brandRepository) GetBrandById(id string) chan RepositoryResult[models.Brand] {
	result := make(chan RepositoryResult[models.Brand])
	go func() {
		brand := models.Brand{}
		err := b.db.Find(&brand).Where("id = ?", id).Error
		if err == gorm.ErrRecordNotFound {
			result <- RepositoryResult[models.Brand]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Brand not found",
			}
			return
		} else if err != nil {
			result <- RepositoryResult[models.Brand]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusInternalServerError,
				Message:    "Error: " + err.Error(),
			}
			return
		} else if brand.Name == "" {
			result <- RepositoryResult[models.Brand]{
				Data:       brand,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Brand not found",
			}
			return
		}
		result <- RepositoryResult[models.Brand]{
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

func (b *brandRepository) CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand] {
	result := make(chan RepositoryResult[models.Brand])
	userName := helper.ExtractUserIDFromToken(tokenString)
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
		result <- RepositoryResult[models.Brand]{
			Data:       brand,
			Error:      nil,
			StatusCode: http.StatusCreated,
			Message:    "Success Create Brand",
		}
	}()
	return result
}
