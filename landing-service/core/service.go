package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	pb "landing-service/proto"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type LandingService interface {
	kldgaInquire(request KldgaInquireRequest) (string, error)
	kldgaCompetition(request KldgaCompetitionRequest) (string, error)
	sendAuthCode(phone string) (string, error)
	verifyAuthCode(phone string, code string) (string, error)
	adapfitInquire(request AdapfitInquireReqeust) (string, error)
}

type landingService struct {
	emailClient pb.EmailServiceClient
	redisClient *redis.Client
}

func NewLandingService(conn *grpc.ClientConn, redisClient *redis.Client) LandingService {
	emailClient := pb.NewEmailServiceClient(conn)
	return &landingService{
		emailClient: emailClient,
		redisClient: redisClient,
	}
}
func (service *landingService) sendAuthCode(phone string) (string, error) {
	// ✅ 단일 컨텍스트 생성 (타임아웃 5초 설정)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 함수 종료 시 컨텍스트 해제
	// 6자리 랜덤 인증번호 생성
	authCode := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	// Redis에 인증번호 저장 (유효시간: 5분)
	err := service.redisClient.Set(ctx, phone, authCode, time.Minute*5).Err()
	if err != nil {
		log.Printf("Failed to save auth code in Redis: %v", err)
		return "", errors.New("failed to save auth code")
	}

	if err := sendCode(phone, authCode); err != nil {
		return "", err
	}

	log.Printf("Auth code for %s is %s", phone, authCode)

	return "200", nil
}

func (service *landingService) verifyAuthCode(phone string, code string) (string, error) {
	// ✅ 단일 컨텍스트 생성 (타임아웃 5초 설정)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 함수 종료 시 컨텍스트 해제

	// Redis에서 인증번호 조회
	storedCode, err := service.redisClient.Get(ctx, phone).Result()
	if err == redis.Nil {
		return "", errors.New("auth code expired or not found")
	} else if err != nil {
		log.Printf("Failed to get auth code: %v", err)
		return "", errors.New("internal error")
	}

	// 입력된 코드와 비교
	if storedCode == code {
		// 인증 성공 시 Redis에 "인증 완료" 상태 플래그 설정
		err := service.redisClient.Set(ctx, phone+":status", "verified", time.Minute*10).Err()
		if err != nil {
			return "", errors.New("failed to set verified status")
		}
		// 기존 인증번호는 삭제하지 않고 그대로 둠
		return "200", nil
	}
	return "", errors.New("invalid auth code")
}

func (service *landingService) kldgaCompetition(request KldgaCompetitionRequest) (string, error) {
	// ✅ 단일 컨텍스트 생성 (타임아웃 5초 설정)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 함수 종료 시 컨텍스트 해제

	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, request.Phone)
	if err != nil || !matched {
		return "", errors.New("invalid phone format, should be 01000000000")
	}
	// Redis에서 인증 상태 확인
	status, err := service.redisClient.Get(ctx, request.Phone+":status").Result()
	if err == redis.Nil || status != "verified" {
		return "", errors.New("phone number not verified")
	} else if err != nil {
		log.Printf("Failed to check verification status: %v", err)
		return "", errors.New("internal error")
	}

	reponse, err := service.emailClient.KldgaSendCompetitionEmail(ctx, &pb.KldgaCompetitionRequest{
		Name:   request.Name,
		Gender: int32(request.Gender),
		League: request.League,
		Career: request.Career,
		Phone:  request.Phone,
		Memo:   request.Memo,
	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("failed to send")
	}

	// 인증 상태 플래그 삭제
	service.redisClient.Del(context.Background(), request.Phone+":status")

	log.Printf("send email: %v", reponse)
	return "200", nil
}

func (service *landingService) kldgaInquire(request KldgaInquireRequest) (string, error) {
	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, request.Phone)
	if err != nil || !matched {
		return "", errors.New("invalid phone format, should be 01000000000")
	}
	reponse, err := service.emailClient.KldgaSendEmail(context.Background(), &pb.KldgaEmailRequest{
		Email:   request.Email,   // 문의한 사람의 이메일
		Name:    request.Name,    // 이름
		Content: request.Content, // 문의 내용
		Phone:   request.Phone,   // 휴대번호
	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("failed to send")
	}

	log.Printf("send email: %v", reponse)
	return "200", nil
}

func (service *landingService) adapfitInquire(request AdapfitInquireReqeust) (string, error) {
	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, request.Phone)
	if err != nil || !matched {
		return "", errors.New("invalid phone format, should be 01000000000")
	}
	reponse, err := service.emailClient.AdapfitInquire(context.Background(), &pb.AdapfitReqeust{
		Name:    request.Name, // 이름
		Class:   request.Class,
		Phone:   request.Phone, // 휴대번호
		Email:   request.Email, // 문의한 사람의 이메일
		Purpose: request.Purpose,
		Career:  request.Career,
		Content: request.Content, // 문의 내용

	})

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("failed to send")
	}

	log.Printf("send email: %v", reponse)
	return "200", nil
}
