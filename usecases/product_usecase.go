package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type ProductUseCaseInterface interface {
	CreateProduct(tokenString string, pr models.ProductRequest) repositories.RepositoryResult[any]
	GetProductById(id string) repositories.RepositoryResult[any]
	GetProducts(pr models.ProductRequest) repositories.RepositoryResult[any]
}

type productUseCase struct {
	r repositories.ProductRepositoryInterface
}

func InitProductUseCase(r repositories.ProductRepositoryInterface) ProductUseCaseInterface {
	return &productUseCase{
		r: r,
	}
}

func (u *productUseCase) GetProducts(pr models.ProductRequest) repositories.RepositoryResult[any] {
	return <-u.r.GetProducts(pr)
}

func (u *productUseCase) GetProductById(id string) repositories.RepositoryResult[any] {
	return <-u.r.GetProductById(id)
}

func (u *productUseCase) CreateProduct(tokenString string, pr models.ProductRequest) repositories.RepositoryResult[any] {
	return <-u.r.CreateProduct(pr, tokenString)
}
