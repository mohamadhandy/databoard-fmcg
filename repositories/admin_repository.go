package repositories

import (
	"klikdaily-databoard/models"
	"net/http"

	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	CreateAdmin(admin models.AdminRequest) chan RepositoryResult[models.Admin]
	GetAdmins(admin models.AdminRequest) chan RepositoryResult[[]models.Admin]
	GetAdminById(id string) chan RepositoryResult[models.Admin]
}

type adminRepository struct {
	db *gorm.DB
}

func InitAdminRepository(db *gorm.DB) AdminRepositoryInterface {
	return &adminRepository{
		db: db,
	}
}

func (r *adminRepository) CreateAdmin(admin models.AdminRequest) chan RepositoryResult[models.Admin] {
	result := make(chan RepositoryResult[models.Admin])
	go func() {
		adminReq := admin.ForCreation()
		r.db.Create(&adminReq)
		result <- RepositoryResult[models.Admin]{
			Data:       adminReq,
			Error:      nil,
			StatusCode: http.StatusCreated,
			Message:    "Success Create Admin",
		}
	}()
	return result
}

func (r *adminRepository) GetAdmins(admin models.AdminRequest) chan RepositoryResult[[]models.Admin] {
	result := make(chan RepositoryResult[[]models.Admin])
	go func() {
		var admins []models.Admin
		adminCount := int64(0)
		page, limit := admin.ForList()
		r.db.Count(&adminCount)
		offset := (page - 1) * limit
		r.db.Offset(int(offset)).Limit(int(limit)).Find(&admins)
		result <- RepositoryResult[[]models.Admin]{
			Data:       admins,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success Get Admins",
		}
	}()
	return result
}

func (r *adminRepository) GetAdminById(id string) chan RepositoryResult[models.Admin] {
	result := make(chan RepositoryResult[models.Admin])
	go func() {
		var admin models.Admin
		r.db.Where("id = ?", id).First(&admin)
		if admin.ID == "" {
			result <- RepositoryResult[models.Admin]{
				Data:       admin,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Admin Not Found",
			}
			return
		}
		result <- RepositoryResult[models.Admin]{
			Data:       admin,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success Get Admin By Id",
		}
	}()
	return result
}
