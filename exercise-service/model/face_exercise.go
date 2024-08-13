package model

import (
	"gorm.io/gorm"
)

type FaceExercise struct {
	gorm.Model
	Id      uint
	Video   Video `gorm:"foreignKey:VideoID"`
	VideoID string
	Type    uint
	Title   string
}
