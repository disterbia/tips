package model

import (
	"time"

	"gorm.io/gorm"
)

type MedicineTake struct {
	gorm.Model
	User       User `gorm:"foreignKey:Uid"`
	Uid        uint
	Medicine   Medicine `gorm:"foreignKey:MedicineID"`
	MedicineID uint
	Dose       float32
	DateTaken  time.Time `gorm:"type:date"`
	TimeTaken  string
}
