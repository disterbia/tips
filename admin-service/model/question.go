package model

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Name         string
	Phone        string
	HospitalName string
	PossibleTime string
	Email        string
	EntryRoute   string
}
