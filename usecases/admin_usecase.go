package usecases

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AdminUseCaseInterface interface {
	CreateAdmin(admin models.AdminRequest) repositories.RepositoryResult[models.Admin]
	GetAdmins(admin models.AdminRequest) repositories.RepositoryResult[[]models.Admin]
	GetAdminById(id string) repositories.RepositoryResult[models.Admin]
}

type adminUsecase struct {
	AdminRepository repositories.AdminRepositoryInterface
}

func InitAdminUsecase(r repositories.AdminRepositoryInterface) AdminUseCaseInterface {
	return &adminUsecase{
		AdminRepository: r,
	}
}

func (u *adminUsecase) CreateAdmin(admin models.AdminRequest) repositories.RepositoryResult[models.Admin] {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	// return err in here
	// 	return user, helpers.NewUnexpectedError("Bcrypt Error")
	// }
	genUuid := helper.InitUuidHelper().GenerateUUID()
	admin.ID = genUuid
	admin.Password = string(passwordHash)
	result := u.AdminRepository.CreateAdmin(admin)
	adminResult := <-result
	return adminResult
}

func (u *adminUsecase) GetAdmins(admin models.AdminRequest) repositories.RepositoryResult[[]models.Admin] {
	result := u.AdminRepository.GetAdmins(admin)
	adminsResult := <-result
	return adminsResult
}

func (u *adminUsecase) GetAdminById(id string) repositories.RepositoryResult[models.Admin] {
	result := u.AdminRepository.GetAdminById(id)
	adminResult := <-result
	return adminResult
}
