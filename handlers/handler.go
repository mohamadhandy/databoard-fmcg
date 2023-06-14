package handlers

import "klikdaily-databoard/usecases"

func InitVersionOneAdminHandler(u usecases.Repositories) AdminHandlerInterface {
	uc := usecases.GetUseCase(u)
	return InitAdminHandler(uc.AdminUsecase)
}

func InitVersionOneAuthenticationHandler(u usecases.Repositories) AuthenticationHandlerInterface {
	uc := usecases.GetUseCase(u)
	return InitAuthenticationHandler(uc.AuthenticationUseCase)
}

func InitVersionOneBrandHandler(u usecases.Repositories) BrandHandlerInterface {
	uc := usecases.GetUseCase(u)
	return InitBrandHandler(uc.BrandUseCase)
}

func InitVersionOneCategoryHandler(u usecases.Repositories) CategoryHandlerInterface {
	uc := usecases.GetUseCase(u)
	return InitCategoryHandler(uc.CategoryUseCase)
}
