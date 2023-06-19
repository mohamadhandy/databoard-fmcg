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
	CreateAdmin(admin models.AdminRequest) chan RepositoryResult[models.Admin]
	GetAdmins(admin models.AdminRequest) chan RepositoryResult[[]models.Admin]
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
				result <- RepositoryResult[[]models.Admin]{
					Data:       admins,
					Error:      nil,
					StatusCode: http.StatusInternalServerError,
					Message:    "Error",
				}
				return
			}
		} else if err != nil {
			result <- RepositoryResult[[]models.Admin]{
				Data:       admins,
				Error:      nil,
				StatusCode: http.StatusInternalServerError,
				Message:    "Error",
			}
			return
		} else {
			err = json.Unmarshal([]byte(value), &admins)
			if err != nil {
				result <- RepositoryResult[[]models.Admin]{
					Data:       admins,
					Error:      nil,
					StatusCode: http.StatusInternalServerError,
					Message:    "Error",
				}
				return
			}
		}

		result <- RepositoryResult[[]models.Admin]{
			Data:       admins,
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
		result <- RepositoryResult[any]{
			Data:       admin,
			Error:      nil,
			StatusCode: http.StatusOK,
			Message:    "Success Get Admin By Id",
		}
	}()
	return result
}
