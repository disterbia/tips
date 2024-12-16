package core

import (
	"context"
	"errors"
	"regexp"

	"log"

	pb "landing-service/proto"

	"google.golang.org/grpc"
)

type LandingService interface {
	kldgaInquire(request KldgaRequest) (string, error)
}

type landingService struct {
	emailClient pb.EmailServiceClient
}

func NewLandingService(conn *grpc.ClientConn) LandingService {
	emailClient := pb.NewEmailServiceClient(conn)
	return &landingService{
		emailClient: emailClient,
	}
}

func (service *landingService) kldgaInquire(request KldgaRequest) (string, error) {
	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, request.Phone)
	if err != nil || !matched {
		return "", errors.New("invalid phone format, should be 01000000000")
	}
	reponse, err := service.emailClient.KldgaSendEmail(context.Background(), &pb.KldgaEmailRequest{
		Email:   request.Email,   // 받는 사람의 이메일
		Name:    request.Name,    // 이름
		Content: request.Content, // 문의 내용
		Phone:   request.Phone,   // 휴대번호
	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("failed to send email")
	}

	log.Printf("send email: %v", reponse)
	return "200", nil
}
