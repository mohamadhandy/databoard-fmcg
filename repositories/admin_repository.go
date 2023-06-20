package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"klikdaily-databoard/models"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdminRepositoryInterface interface {
	CreateAdmin(admin models.AdminRequest) chan RepositoryResult[any]
	GetAdmins(admin models.AdminRequest) chan RepositoryResult[any]
	GetAdminById(id string) chan RepositoryResult[any]
}

type adminRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func InitAdminRepository(db *gorm.DB, rdb *redis.Client) AdminRepositoryInterface {
	return &adminRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *adminRepository) CreateAdmin(admin models.AdminRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		adminReq := admin.ForCreation()
		if err := r.db.Create(&adminReq).Error; err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[any]{
			Data:       adminReq,
			Error:      nil,
			StatusCode: http.StatusCreated,
			Message:    "Success Create Admin",
		}
	}()
	return result
}

func (r *adminRepository) GetAdmins(admin models.AdminRequest) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	ctx := context.Background()
	go func() {
		var admins []models.Admin
		adminCount := int64(0)
		page, limit := admin.ForList()
		value, err := r.rdb.Get(ctx, fmt.Sprintf("page_%v_limit_%v", page, limit)).Result()
		if err == redis.Nil {
			r.db.Count(&adminCount)
			offset := (page - 1) * limit
			r.db.Offset(int(offset)).Limit(int(limit)).Find(&admins)
			valueMarshal, _ := json.Marshal(admins)
			err = r.rdb.Set(ctx, fmt.Sprintf("page_%v_limit_%v", page, limit), string(valueMarshal), 1*time.Minute).Err()
			if err != nil {
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      nil,
					StatusCode: http.StatusInternalServerError,
					Message:    "Error: " + err.Error(),
				}
				return
			}
		} else if err != nil {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      nil,
				StatusCode: http.StatusInternalServerError,
				Message:    "Error: " + err.Error(),
			}
			return
		} else {
			err = json.Unmarshal([]byte(value), &admins)
			if err != nil {
				result <- RepositoryResult[any]{
					Data:       nil,
					Error:      nil,
					StatusCode: http.StatusInternalServerError,
					Message:    "Error: " + err.Error(),
				}
				return
			}
		}
		adminResponses := make([]models.AdminResponse, len(admins))
		for i, v := range admins {
			adminResponses[i] = models.AdminResponse{
				Email:       v.Email,
				Name:        v.Name,
				ID:          v.ID,
				Phonenumber: v.PhoneNumber,
				Status:      v.Status,
				Role:        v.Role,
			}
		}
		result <- RepositoryResult[any]{
			Data:       adminResponses,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success Get Admins",
		}
	}()
	return result
}

func (r *adminRepository) GetAdminById(id string) chan RepositoryResult[any] {
	result := make(chan RepositoryResult[any])
	go func() {
		var admin models.Admin
		r.db.Where("id = ?", id).First(&admin)
		if admin.ID == "" {
			result <- RepositoryResult[any]{
				Data:       nil,
				Error:      nil,
				StatusCode: http.StatusNotFound,
				Message:    "Admin Not Found",
			}
			return
		}
		adminResponse := models.AdminResponse{
			Name:        admin.Name,
			Email:       admin.Email,
			ID:          admin.ID,
			Phonenumber: admin.PhoneNumber,
			Status:      admin.Status,
			Role:        admin.Role,
		}
		result <- RepositoryResult[any]{
			Data:       adminResponse,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success Get Admin By Id",
		}
	}()
	return result
}
