package model

import (
	"gorm.io/gorm"
)

type LinkedEmail struct {
	gorm.Model
	Email   string
	User    User `gorm:"foreignKey:Uid"`
	Uid     uint
	SnsType uint
}
