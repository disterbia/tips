package model

import (
	"gorm.io/gorm"
)

type InquireReply struct {
	gorm.Model
	Id        uint
	User      User `gorm:"foreignKey:Uid"`
	Uid       uint
	Inquire   Inquire `gorm:"foreignKey:InquireID"`
	InquireID uint
	ReplyType bool
	Content   string
}
