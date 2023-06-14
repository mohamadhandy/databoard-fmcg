package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type CategoryUseCaseInterface interface {
	CreateCategory(tokenString string, cr models.CategoryRequest) repositories.RepositoryResult[models.Category]
}

type categoryUseCase struct {
	r repositories.CategoryRepositoryInterface
}

func InitCategoryUseCase(r repositories.CategoryRepositoryInterface) CategoryUseCaseInterface {
	return &categoryUseCase{
		r,
	}
}

func (c *categoryUseCase) CreateCategory(tokenString string, cr models.CategoryRequest) repositories.RepositoryResult[models.Category] {
	return <-c.r.CreateCategory(tokenString, cr)
}
