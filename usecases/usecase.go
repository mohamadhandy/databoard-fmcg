package usecases

import (
	"klikdaily-databoard/repositories"

	"klikdaily-databoard/config"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	AdminRepository          repositories.AdminRepositoryInterface
	AuthenticationRepository repositories.AuthenticationRepositoryInterface
	BrandRepository          repositories.BrandRepositoryInterface
	CategoryRepository       repositories.CategoryRepositoryInterface
	ProductRepository        repositories.ProductRepositoryInterface
}

type Usecases struct {
	AdminUsecase          AdminUseCaseInterface
	AuthenticationUseCase AuthenticationUseCaseInterface
	BrandUseCase          BrandUsecaseInterface
	CategoryUseCase       CategoryUseCaseInterface
	ProductUseCase        ProductUseCaseInterface
}

var useCaseInstance Usecases

func InitRepository(db *gorm.DB, rdb *redis.Client, es *elasticsearch.Client, mb *config.MessageBroker) Repositories {
	return Repositories{
		AdminRepository:          repositories.InitAdminRepository(db, rdb),
		AuthenticationRepository: repositories.InitAuthenticationRepository(db),
		BrandRepository:          repositories.InitBrandRepository(db),
		CategoryRepository:       repositories.InitCategoryRepository(db),
		ProductRepository:        repositories.InitProductRepository(db, es, mb),
	}
}

func GetUseCase(r Repositories) *Usecases {
	if (Usecases{}) == useCaseInstance {
		useCaseInstance = Usecases{
			AdminUsecase:          InitAdminUsecase(r.AdminRepository),
			AuthenticationUseCase: InitAuthenticationUseCase(r.AuthenticationRepository),
			BrandUseCase:          InitBrandUseCase(r.BrandRepository),
			CategoryUseCase:       InitCategoryUseCase(r.CategoryRepository),
			ProductUseCase:        InitProductUseCase(r.ProductRepository),
		}
	}
	return &useCaseInstance
}
