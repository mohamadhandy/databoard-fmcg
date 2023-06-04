package usecases

import (
	"klikdaily-databoard/repositories"

	"gorm.io/gorm"
)

type Repositories struct {
	AdminRepository          repositories.AdminRepositoryInterface
	AuthenticationRepository repositories.AuthenticationRepositoryInterface
}

type Usecases struct {
	AdminUsecase          AdminUseCaseInterface
	AuthenticationUseCase AuthenticationUseCaseInterface
}

var useCaseInstance Usecases

func InitRepository(db *gorm.DB) Repositories {
	return Repositories{
		AdminRepository:          repositories.InitAdminRepository(db),
		AuthenticationRepository: repositories.InitAuthenticationRepository(db),
	}
}

func GetUseCase(r Repositories) *Usecases {
	if (Usecases{}) == useCaseInstance {
		useCaseInstance = Usecases{
			AdminUsecase:          InitAdminUsecase(r.AdminRepository),
			AuthenticationUseCase: InitAuthenticationUseCase(r.AuthenticationRepository),
		}
	}
	return &useCaseInstance
}
