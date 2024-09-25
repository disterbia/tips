package model

import (
	"gorm.io/gorm"
)

type TapBlinkScore struct {
	gorm.Model
	User          User `gorm:"foreignKey:Uid"`
	Uid           uint
	ScoreType     uint //tap or blink
	SuccessCount  uint
	ErrorCount    uint
	ReactionSpeed float64
}
