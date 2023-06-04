package usecases

import (
	"klikdaily-databoard/repositories"

	"gorm.io/gorm"
)

type Repositories struct {
	AdminRepository repositories.AdminRepositoryInterface
}

type Usecases struct {
	AdminUsecase AdminUseCaseInterface
}

var useCaseInstance Usecases

func InitRepository(db *gorm.DB) Repositories {
	return Repositories{
		AdminRepository: repositories.InitAdminRepository(db),
	}
}

func GetUseCase(r Repositories) *Usecases {
	if (Usecases{}) == useCaseInstance {
		useCaseInstance = Usecases{
			AdminUsecase: InitAdminUsecase(r.AdminRepository),
		}
	}
	return &useCaseInstance
}
