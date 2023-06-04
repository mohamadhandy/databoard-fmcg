package handlers

import (
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"
	"strconv"

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
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "2")

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)
	adminRequest := models.AdminRequest{
		Page:  uint(p),
		Limit: uint(l),
	}
	adminsResult := h.AdminUsecase.GetAdmins(adminRequest)
	c.JSON(adminsResult.StatusCode, adminsResult)
}

func (h *adminHandler) GetAdminById(c *gin.Context) {
	id := c.Param("id")
	adminResult := h.AdminUsecase.GetAdminById(id)
	c.JSON(adminResult.StatusCode, adminResult)
}
