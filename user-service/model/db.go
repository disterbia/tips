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

	// 1. 기존 유니크 인덱스가 있으면 삭제
	//    인덱스 이름은 DB에서 직접 확인하고 수정 가능
	db.Exec("DROP INDEX IF EXISTS uni_users_email;")

	// 2. 조건부 인덱스를 생성: DeletedAt이 NULL인 경우에만 유니크 적용
	db.Exec("CREATE UNIQUE INDEX unique_email_on_users ON users (email) WHERE deleted_at IS NULL;")

	db.AutoMigrate(&AppVersion{}, &AuthCode{}, &Image{}, &LinkedEmail{}, &User{}, &VerifiedTarget{}, &Police{})
	return db, nil
}
