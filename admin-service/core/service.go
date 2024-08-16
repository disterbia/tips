// /user-service/service/service.go

package core

import (
	"errors"
	"strings"

	"admin-service/model"

	pb "admin-service/proto"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type AdminService interface {
	login(request LoginRequest) (string, error)
}

type adminService struct {
	db          *gorm.DB
	emailClient pb.EmailServiceClient
}

func NewAdminService(db *gorm.DB, conn *grpc.ClientConn) AdminService {
	emailClient := pb.NewEmailServiceClient(conn)
	return &adminService{
		db:          db,
		emailClient: emailClient,
	}
}

func (service *adminService) login(request LoginRequest) (string, error) {
	var u model.User
	password := strings.TrimSpace(request.Password)

	if password == "" {
		return "", errors.New("empty")
	}

	// 이메일로 사용자 조회
	if err := service.db.Where("email = ?", request.Email).First(&u).Error; err != nil {
		return "", errors.New("-2")
	}

	// 비밀번호 비교
	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(request.Password)); err != nil {
		return "", errors.New("-2")
	}

	if !*u.IsApproval {
		return "", errors.New("-1")
	}

	// 새로운 JWT 토큰 생성
	tokenString, err := generateJWT(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
