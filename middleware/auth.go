package middleware

import (
	"klikdaily-databoard/helper"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			helper.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "unauthorized error parse claims!")
			return
		}
		if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized err")
		}
	}
}
