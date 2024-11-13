package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model    // ID, CreatedAt, UpdatedAt, DeletedAt 필드를 자동으로 추가
	Name          string
	Email         *string
	DeviceID      string
	FCMToken      string
	SnsType       uint
	Phone         string
	Gender        bool
	Birthday      time.Time `gorm:"type:date"`
	UserType      uint
	ProfileImages []Image       `gorm:"foreignKey:Uid"`
	LinkedEmails  []LinkedEmail `gorm:"foreignkey:Uid"`
	UserPolices   []UserPolice  `gorm:"foreignkey:Uid"`
}
