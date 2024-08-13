package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model   // ID, CreatedAt, UpdatedAt, DeletedAt 필드를 자동으로 추가
	ProjectName  string
	Name         string
	Duration     uint
	ProjectId    string
	VideoId      string
	ThumbnailUrl string
}
