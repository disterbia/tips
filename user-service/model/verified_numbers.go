package model

import (
	"gorm.io/gorm"
)

type VerifiedNumbers struct {
	gorm.Model
	Id    uint
	Phone string
}
