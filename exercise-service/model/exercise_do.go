package model

import (
	"time"

	"gorm.io/gorm"
)

type ExerciseTake struct {
	gorm.Model
	User       User `gorm:"foreignKey:Uid"`
	Uid        uint
	Exercise   Exercise `gorm:"foreignKey:ExerciseID"`
	ExerciseID uint
	DateTaken  time.Time `gorm:"type:date"`
	TimeTaken  string
}
