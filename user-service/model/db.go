package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 데이터베이스 연결 초기화
func NewDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	// lock_timeout 설정
	if err = db.Exec("SET lock_timeout = '1s'").Error; err != nil {
		return nil, err
	}

	db.AutoMigrate(&AppVersion{}, &AuthCode{}, &Image{}, &LinkedEmail{}, &User{}, &VerifiedTarget{}, &Polices{})
	return db, nil
}
