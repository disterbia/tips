package core

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

var jwtSecretKey = []byte("adapfit_mark")

func verifyJWT(c *fiber.Ctx) (uint, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return 0, errors.New("authorization header is required")
	}

	// 'Bearer ' 접두사 제거
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
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

func validateExercise(medicine ExerciseRequest) error {
	if medicine.StartAt != "" {
		if err := validateDate(medicine.StartAt); err != nil {
			return err
		}
	}
	if medicine.EndAt != "" {
		if err := validateDate(medicine.EndAt); err != nil {
			return err
		}
	}

	return nil
}

func validateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.New("invalid date format, should be YYYY-MM-DD")
	}
	return nil
}

func validateTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return errors.New("invalid time format, should be HH:MM")
	}
	return nil
}

// 배열에 특정 값이 포함되어 있는지 확인하는 함수
func contains(arr []int64, target int64) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}
