package handlers

import (
	"klikdaily-databoard/helper"
	"klikdaily-databoard/middleware"
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationHandlerInterface interface {
	Login(c *gin.Context)
}

type AuthenticationHandler struct {
	auth usecases.AuthenticationUseCaseInterface
}

func InitAuthenticationHandler(auth usecases.AuthenticationUseCaseInterface) AuthenticationHandlerInterface {
	return &AuthenticationHandler{
		auth,
	}
}

func (h *AuthenticationHandler) Login(c *gin.Context) {
	email := c.Request.FormValue("email")
	pass := c.Request.FormValue("password")
	reqModel := models.RequestableAuthentication{
		Email:    email,
		Password: pass,
	}
	result := h.auth.BeginSession(&reqModel)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, result)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.MyCustomClaims{
		Email: result.Data.Email,
		Name:  result.Data.Name,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		helper.Error(err.Error())
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, tokenString)
}
