package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type ProductUseCaseInterface interface {
	CreateProduct(tokenString string, pr models.ProductRequest) repositories.RepositoryResult[any]
}

type productUseCase struct {
	r repositories.ProductRepositoryInterface
}

func InitProductUseCase(r repositories.ProductRepositoryInterface) ProductUseCaseInterface {
	return &productUseCase{
		r: r,
	}
}

func (u *productUseCase) CreateProduct(tokenString string, pr models.ProductRequest) repositories.RepositoryResult[any] {
	return <-u.r.CreateProduct(pr, tokenString)
}
