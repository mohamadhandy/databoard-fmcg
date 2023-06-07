package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type BrandUsecaseInterface interface {
	CreateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[models.Brand]
}

type brandUsecase struct {
	r repositories.BrandRepositoryInterface
}

func InitBrandUseCase(r repositories.BrandRepositoryInterface) BrandUsecaseInterface {
	return &brandUsecase{
		r,
	}
}

func (b *brandUsecase) CreateBrand(tokenString string, br models.BrandRequest) repositories.RepositoryResult[models.Brand] {
	res := b.r.CreateBrand(tokenString, br)
	brandResult := <-res
	return brandResult
}
