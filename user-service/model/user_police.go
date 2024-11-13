package model

import (
	"gorm.io/gorm"
)

type UserPolice struct {
	gorm.Model
	User     User `gorm:"foreignKey:Uid"`
	Uid      uint
	Police   Police `gorm:"foreignKey:PoliceId"`
	PoliceId uint
}
