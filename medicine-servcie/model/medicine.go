package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Medicine struct {
	gorm.Model
	User         User `gorm:"foreignKey:Uid"`
	Uid          uint
	Name         string
	Weekdays     pq.Int64Array  `gorm:"type:integer[]"`
	Times        pq.StringArray `gorm:"type:text[]"`
	Dose         float32
	MedicineType string
	StartAt      *time.Time `gorm:"type:date"`
	EndAt        *time.Time `gorm:"type:date"`
	MinReserves  *float32
	Remaining    *float32
	UsePrivacy   bool
	IsActive     bool
}
