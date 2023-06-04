package handlers

import "klikdaily-databoard/usecases"

func InitVersionOneAdminHandler(u usecases.Repositories) AdminHandlerInterface {
	uc := usecases.GetUseCase(u)
	return InitAdminHandler(uc.AdminUsecase)
}
