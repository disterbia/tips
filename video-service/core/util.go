package core

import (
	"errors"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtSecretKey = []byte("adapfit_mark")

func verifyJWT(c *fiber.Ctx) (uint, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		log.Println("required")
		return 0, errors.New("authorization header is required")
	}

	// 'Bearer ' 접두사 제거
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("invalid")
		return 0, errors.New("invalid token")
	}

	id := uint((*claims)["id"].(float64))

	return id, nil
}
