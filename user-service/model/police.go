package model

import (
	"gorm.io/gorm"
)

type Polices struct {
	gorm.Model
	Title      string
	Body       string
	PoliceType uint
}
