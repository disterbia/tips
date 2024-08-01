package model

import (
	"gorm.io/gorm"
)

type AuthCode struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt 필드를 자동으로 추가
	Id         uint
	Phone      string
	Code       string
}
