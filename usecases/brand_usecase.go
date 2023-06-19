package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type BrandUsecaseInterface interface {
	CreateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[any]
	GetBrandById(id string) repositories.RepositoryResult[any]
	GetBrands(brandRequest models.BrandRequest) repositories.RepositoryResult[any]
	UpdateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[models.Brand]
}

type brandUsecase struct {
	r repositories.BrandRepositoryInterface
}

func InitBrandUseCase(r repositories.BrandRepositoryInterface) BrandUsecaseInterface {
	return &brandUsecase{
		r,
	}
}

func (b *brandUsecase) GetBrands(brandRequest models.BrandRequest) repositories.RepositoryResult[any] {
	return <-b.r.GetBrands(brandRequest)
}

func (b *brandUsecase) CreateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[any] {
	return <-b.r.CreateBrand(tokenString, br)
}

func (b *brandUsecase) GetBrandById(id string) repositories.RepositoryResult[any] {
	return <-b.r.GetBrandById(id)
}

func (b *brandUsecase) UpdateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[models.Brand] {
	return <-b.r.UpdateBrand(tokenString, br)
}
