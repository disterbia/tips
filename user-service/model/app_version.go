package model

import (
	"gorm.io/gorm"
)

type AppVersion struct {
	gorm.Model
	LatestVersion string
	AndroidLink   string
	IosLink       string
}
