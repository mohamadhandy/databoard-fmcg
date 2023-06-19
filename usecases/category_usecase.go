package usecases

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/repositories"
)

type CategoryUseCaseInterface interface {
	CreateCategory(tokenString string, cr models.CategoryRequest) repositories.RepositoryResult[any]
	GetCategories(tokenString string) repositories.RepositoryResult[any]
}

type categoryUseCase struct {
	r repositories.CategoryRepositoryInterface
}

func InitCategoryUseCase(r repositories.CategoryRepositoryInterface) CategoryUseCaseInterface {
	return &categoryUseCase{
		r,
	}
}

func (c *categoryUseCase) CreateCategory(tokenString string, cr models.CategoryRequest) repositories.RepositoryResult[any] {
	return <-c.r.CreateCategory(tokenString, cr)
}

func (c *categoryUseCase) GetCategories(tokenString string) repositories.RepositoryResult[any] {
	return <-c.r.GetCategories(tokenString)
}
