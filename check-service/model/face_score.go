package model

import (
	"gorm.io/gorm"
)

type FaceScore struct {
	gorm.Model
	User     User `gorm:"foreignKey:Uid"`
	Uid      uint
	FaceType uint
	FaceLine uint
	Sd       float64 //smoothing in distance
}
