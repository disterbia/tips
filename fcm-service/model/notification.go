package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	User     User `gorm:"foreignKey:Uid"`
	Uid      uint
	ParentId uint
	Type     uint
	Body     string
	StartAt  *time.Time `gorm:"type:date"`
	EndAt    *time.Time `gorm:"type:date"`
	Time     string
	Weekdays pq.Int64Array `gorm:"type:integer[]"`
}
