package repositories

import (
	"klikdaily-databoard/models"
	"net/http"

	"gorm.io/gorm"
)

type BrandRepositoryInterface interface {
	CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand]
}

type brandRepository struct {
	db *gorm.DB
}

func InitBrandRepository(db *gorm.DB) BrandRepositoryInterface {
	return &brandRepository{
		db,
	}
}

func (b *brandRepository) CreateBrand(tokenString string, br models.BrandRequest) chan RepositoryResult[models.Brand] {
	result := make(chan RepositoryResult[models.Brand])
	// getTokenString := helper.ExtractUserIDFromToken(tokenString)
	// fmt.Println("test token string", tokenString)
	// fmt.Println(getTokenString)
	// fmt.Println("test1234")
	go func() {
		// validate brandRequest is created_by / updated_by superadmin role or adminGudang role or warehouseLeader role ?
		brand := models.Brand{
			Name: br.Name,
			ID:   "001",
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
