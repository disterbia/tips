package model

import (
	"gorm.io/gorm"
)

type AppVersion struct {
	gorm.Model
	Id            uint
	LatestVersion string
	AndroidLink   string
	IosLink       string
}
