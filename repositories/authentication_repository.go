package repositories

import (
	"klikdaily-databoard/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationRepositoryInterface interface {
	BeginSession(request models.RequestableAuthenticationInterface) chan RepositoryResult[models.Admin]
}

type AuthenticationRepository struct {
	db *gorm.DB
}

func InitAuthenticationRepository(db *gorm.DB) AuthenticationRepositoryInterface {
	return &AuthenticationRepository{db}
}

func (a *AuthenticationRepository) BeginSession(request models.RequestableAuthenticationInterface) chan RepositoryResult[models.Admin] {
	result := make(chan RepositoryResult[models.Admin])
	go func() {
		admin := models.Admin{}
		email, pass := request.ForAuthentication()
		err := a.db.Where(&models.Admin{
			Email: email,
		}).First(&admin).Error
		if err == gorm.ErrRecordNotFound {
			// isi data nanti
			result <- RepositoryResult[models.Admin]{
				Data:       admin,
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
			return
		} else if err != nil {
			// isi data nanti
			result <- RepositoryResult[models.Admin]{
				Data:       admin,
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(pass)); err != nil {
			result <- RepositoryResult[models.Admin]{
				Data:       admin,
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
			return
		}
		result <- RepositoryResult[models.Admin]{
			Data:       admin,
			Error:      nil,
			StatusCode: http.StatusOK,
		}
	}()
	return result
}
