package model

import (
	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model
	Title  string
	Body   string
	IsLast bool
}
