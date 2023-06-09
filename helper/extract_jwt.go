package helper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractUserIDFromToken(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		// Verify the token using the appropriate signing key or method for your application
		// For example, if you're using HMAC signing method:
		signingKey := []byte(os.Getenv("SECRET_KEY"))
		return signingKey, nil
	})
	fmt.Println("token: " + token.Raw)
	if err != nil {
		// Handle token parsing or verification errors
		return "Parsing error: " + err.Error()
	}

	// Extract the user ID from the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["name"].(string); ok {
			return userID
		}
	}

	return "Token not valid"
}
