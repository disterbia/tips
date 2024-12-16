// /email-service/service/email-service.go

package service

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	pb "email-service/proto"
)

type EmailServer struct {
	pb.UnimplementedEmailServiceServer
}

func (s *EmailServer) SendEmail(ctx context.Context, req *pb.EmailRequest) (*pb.EmailResponse, error) {
	// SMTP 설정
	email := os.Getenv("WELLKINSON_SMTP_EMAIL")
	password := os.Getenv("WELLKINSON_SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 인증 정보
	auth := smtp.PlainAuth("", email, password, smtpHost)

	// 이메일 본문 구성
	body := fmt.Sprintf(
		"<h2>작성자: </h2><span>%s</span><br>"+
			"<h2>날짜: </h2><span>%s</span><br>"+
			"<h2>제목: </h2><span>%s</span><br>"+
			"<h2>내용: </h2><span>%s</span><br>"+
			"<h2>답변: </h2><span>%s</span><br>"+
			"<h2>답변 날짜: </h2><span>%s</span><br>",
		req.Email, req.CreatedAt, req.Title, req.Content, req.ReplyContent, req.ReplyCreatedAt)

	// 이메일 메시지 설정
	msg := []byte("To: " + req.Email + "\r\n" +
		"Subject: 문의 주신 " + req.Title + "에 답변드립니다.\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	// 이메일 전송
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, []string{req.Email}, msg)
	if err != nil {
		return nil, err
	}

	return &pb.EmailResponse{Status: "Success"}, nil
}

func (s *EmailServer) SendCodeEmail(ctx context.Context, req *pb.EmailCodeRequest) (*pb.EmailResponse, error) {
	// SMTP 설정
	email := os.Getenv("WELLKINSON_SMTP_EMAIL")
	password := os.Getenv("WELLKINSON_SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 인증 정보
	auth := smtp.PlainAuth("", email, password, smtpHost)

	// 이메일 본문 구성
	body := fmt.Sprintf(
		"<h2>인증번호: </h2><span>%s</span><br>"+
			"<h3>위의 인증번호를 입력해주세요.</h3>",
		req.Code)

	// 이메일 메시지 설정
	msg := []byte("To: " + req.Email + "\r\n" +
		"Subject:DAF 회원가입 인증번호 입니다. \r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	// 이메일 전송
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, []string{req.Email}, msg)
	if err != nil {
		return nil, err
	}

	return &pb.EmailResponse{Status: "Success"}, nil
}

func (s *EmailServer) KldgaSendEmail(ctx context.Context, req *pb.KldgaEmailRequest) (*pb.EmailResponse, error) {
	// SMTP 설정
	email := os.Getenv("WELLKINSON_SMTP_EMAIL")
	password := os.Getenv("WELLKINSON_SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 인증 정보
	auth := smtp.PlainAuth("", email, password, smtpHost)

	// 이메일 본문 구성
	body := fmt.Sprintf(
		"<h1>아래와 같은 문의가 등록되었습니다.<h1><br>"+
			"<h2>이름: </h2><span>%s</span><br>"+
			"<h2>이메일: </h2><span>%s</span><br>"+
			"<h2>휴대번호: </h2><span>%s</span><br>"+
			"<h2>내용: </h2><span>%s</span><br>",
		req.Name, req.Email, req.Phone, req.Content)

	// 이메일 메시지 설정
	msg := []byte("To: disterbia@naver.com\r\n" +
		"Subject: kldga 문의등록 알림\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	// 이메일 전송
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, []string{"disterbia@naver.com"}, msg)
	if err != nil {
		return nil, err
	}

	return &pb.EmailResponse{Status: "Success"}, nil
}
