package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateLocalJWT() string {
	claims := jwt.MapClaims{
		"sub":            "test-user-id",
		"email":          "test@example.com",
		"cognito:groups": []string{"Users"},
		"exp":            time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("LOCAL_JWT_SECRET")))
	if err != nil {
		log.Fatal("failed to create JWT:", err)
	}
	return tokenString
}
