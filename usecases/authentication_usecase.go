package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type AuthenticationUseCaseInterface interface {
	BeginSession(models.RequestableAuthenticationInterface) repositories.RepositoryResult[models.Admin]
}

type AuthenticationUseCase struct {
	r repositories.AuthenticationRepositoryInterface
}

func InitAuthenticationUseCase(r repositories.AuthenticationRepositoryInterface) AuthenticationUseCaseInterface {
	return &AuthenticationUseCase{
		r,
	}
}

func (u *AuthenticationUseCase) BeginSession(req models.RequestableAuthenticationInterface) repositories.RepositoryResult[models.Admin] {
	result := u.r.BeginSession(req)
	res := <-result
	return res
}
