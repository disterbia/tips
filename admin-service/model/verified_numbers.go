package model

import (
	"gorm.io/gorm"
)

type VerifiedTarget struct {
	gorm.Model
	Target string
}
