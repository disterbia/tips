package model

import (
	"gorm.io/gorm"
)

type Hospital struct {
	gorm.Model
	Name       string
	Number     string
	Address    string
	RegionCode string `gorm:"index"`
}
