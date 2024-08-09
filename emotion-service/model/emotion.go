package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Emotion struct {
	gorm.Model
	User       User `gorm:"foreignKey:Uid"`
	Uid        uint
	Emotion    uint
	Symptoms   pq.Int64Array `gorm:"type:integer[]"`
	Memo       string
	TargetDate time.Time `gorm:"type:date"`
}
