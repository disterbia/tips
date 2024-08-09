package model

import (
	"gorm.io/gorm"
)

type MedicineSearch struct {
	gorm.Model
	Name string
}
