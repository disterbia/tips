package model

import (
	"gorm.io/gorm"
)

type Inquire struct {
	gorm.Model
	Id      uint
	User    User `gorm:"foreignKey:Uid"`
	Uid     uint
	Email   string
	Title   string
	Content string
	Replies []InquireReply `gorm:"foreignKey:InquireID"`
}
