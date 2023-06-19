package usecases

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AdminUseCaseInterface interface {
	CreateAdmin(admin models.AdminRequest) repositories.RepositoryResult[any]
	GetAdmins(admin models.AdminRequest) repositories.RepositoryResult[any]
	GetAdminById(id string) repositories.RepositoryResult[any]
}

type adminUsecase struct {
	AdminRepository repositories.AdminRepositoryInterface
}

func InitAdminUsecase(r repositories.AdminRepositoryInterface) AdminUseCaseInterface {
	return &adminUsecase{
		AdminRepository: r,
	}
}

func (u *adminUsecase) CreateAdmin(admin models.AdminRequest) repositories.RepositoryResult[any] {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		// return err
		return repositories.RepositoryResult[any]{
			Data:       nil,
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	genUuid := helper.InitUuidHelper().GenerateUUID()
	admin.ID = genUuid
	admin.Password = string(passwordHash)
	return <-u.AdminRepository.CreateAdmin(admin)
}

func (u *adminUsecase) GetAdmins(admin models.AdminRequest) repositories.RepositoryResult[any] {
	return <-u.AdminRepository.GetAdmins(admin)
}

func (u *adminUsecase) GetAdminById(id string) repositories.RepositoryResult[any] {
	return <-u.AdminRepository.GetAdminById(id)
}
