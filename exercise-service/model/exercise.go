package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	User     User `gorm:"foreignKey:Uid"`
	Uid      uint
	Name     string
	Weekdays pq.Int64Array  `gorm:"type:integer[]"`
	Times    pq.StringArray `gorm:"type:text[]"`
	StartAt  *time.Time     `gorm:"type:date"`
	EndAt    *time.Time     `gorm:"type:date"`
	IsActive bool
}
