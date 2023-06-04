package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"

	"github.com/gin-gonic/gin"
)

type AdminHandlerInterface interface {
	CreateAdmin(c *gin.Context)
	GetAdmins(c *gin.Context)
	GetAdminById(c *gin.Context)
}

type adminHandler struct {
	AdminUsecase usecases.AdminUseCaseInterface
}

func InitAdminHandler(u usecases.AdminUseCaseInterface) AdminHandlerInterface {
	return &adminHandler{
		AdminUsecase: u,
	}
}

func (h *adminHandler) CreateAdmin(c *gin.Context) {
	adminRequest := models.AdminRequest{}
	c.BindJSON(&adminRequest)
	adminResult := h.AdminUsecase.CreateAdmin(adminRequest)
	c.JSON(adminResult.StatusCode, adminResult)
}

func (h *adminHandler) GetAdmins(c *gin.Context) {
	adminsResult := h.AdminUsecase.GetAdmins()
	c.JSON(adminsResult.StatusCode, adminsResult)
}

func (h *adminHandler) GetAdminById(c *gin.Context) {
	id := c.Param("id")
	adminResult := h.AdminUsecase.GetAdminById(id)
	c.JSON(adminResult.StatusCode, adminResult)
}
