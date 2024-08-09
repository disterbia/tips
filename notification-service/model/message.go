package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	User     User `gorm:"foreignKey:Uid"`
	Uid      uint
	Type     uint
	Body     string
	ParentID uint
	IsRead   bool
}
