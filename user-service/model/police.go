package model

import (
	"gorm.io/gorm"
)

type Police struct {
	gorm.Model
	Title      string
	Body       string
	PoliceType uint
	IsLast     bool
}
