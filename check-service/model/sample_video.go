package model

import (
	"gorm.io/gorm"
)

type SampleVideo struct {
	gorm.Model
	Category  uint
	Titile    string
	VideoType uint
	VideoId   string
}
