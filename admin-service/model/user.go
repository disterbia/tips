package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt 필드를 자동으로 추가
	Name       string
	Email      *string `gorm:"unique"`
	Password   *string
	DeviceID   string
	FCMToken   string
	SnsType    uint
	Phone      string
	Gender     bool
	Birthday   time.Time `gorm:"type:date"`
	UserType   uint
	Role       Role `gorm:"foreignKey:RoleID"`
	RoleID     uint
	IsApproval *bool
}
