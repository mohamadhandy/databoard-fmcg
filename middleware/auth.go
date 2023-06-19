package middleware

import (
	"fmt"
	"klikdaily-databoard/helper"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			fmt.Println("error: " + err.Error())
			data := helper.Response{
				Data:       nil,
				Error:      err,
				Message:    err.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, data)
			return
		}
		if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, "unauthorized err")
		}
	}
}

func ExtractNameFromToken(tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Println("error extract token name: " + err.Error())
		return err.Error()
	}
	fmt.Println("token test: " + token.Raw)

	// Extract the user ID from the token claims
	if res, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Println("extract token success: " + res.Name)
		return res.Name
	} else {
		fmt.Println("extract token not valid")
		return "Token Not Valid"
	}
}
