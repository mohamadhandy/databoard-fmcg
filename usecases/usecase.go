package usecases

import (
	"klikdaily-databoard/repositories"

	"github.com/redis/go-redis/v9"
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

func InitRepository(db *gorm.DB, rdb *redis.Client) Repositories {
	return Repositories{
		AdminRepository:          repositories.InitAdminRepository(db, rdb),
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
