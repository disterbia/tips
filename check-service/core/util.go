package core

import (
	"errors"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/lib/pq"
)

// JWT secret key
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

func uintSliceToInt64Array(uintSlice []uint) pq.Int64Array {
	int64Array := make(pq.Int64Array, len(uintSlice))
	for i, v := range uintSlice {
		int64Array[i] = int64(v)
	}
	return int64Array
}

func int64ArrayToUintSlice(int64Array pq.Int64Array) []uint {
	uintSlice := make([]uint, len(int64Array))
	for i, v := range int64Array {
		uintSlice[i] = uint(v)
	}
	return uintSlice
}
