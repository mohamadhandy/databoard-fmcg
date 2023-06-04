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
