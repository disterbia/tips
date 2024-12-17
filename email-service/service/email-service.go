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
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f9f9f9; padding: 20px; margin: 0;">
    <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border: 1px solid #ddd; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: #4caf50; color: white; padding: 20px; text-align: center;">
            <h1 style="margin: 0; font-size: 24px;">문의 답변 알림</h1>
        </div>
        
        <!-- Content -->
        <div style="padding: 20px;">
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">작성자:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">날짜:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">제목:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">내용:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">답변:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">답변 날짜:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
        </div>
        
        <!-- Footer -->
        <div style="text-align: center; padding: 10px; background: #f1f1f1; font-size: 12px; color: #666;">
            본 메일은 자동 발송된 메일입니다. 문의사항이 있으시면 관리자에게 연락해주세요.
        </div>
    </div>
</body>
</html>
`, req.Email, req.CreatedAt, req.Title, req.Content, req.ReplyContent, req.ReplyCreatedAt)

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
	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f9f9f9; padding: 20px; margin: 0;">
		<div style="max-width: 600px; margin: 0 auto; background: #ffffff; border: 1px solid #ddd; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); overflow: hidden;">
			<!-- Header -->
			<div style="background: #4caf50; color: white; padding: 20px; text-align: center;">
				<h1 style="margin: 0; font-size: 24px;">인증번호 안내</h1>
			</div>
			
			<!-- Content -->
			<div style="padding: 20px; text-align: center;">
				<p style="font-size: 18px; color: #555;">아래의 인증번호를 입력해주세요.</p>
				<div style="margin: 20px auto; padding: 10px 20px; display: inline-block; background: #f1f1f1; border: 1px solid #ccc; border-radius: 8px;">
					<span style="font-size: 24px; font-weight: bold; color: #4caf50;">%s</span>
				</div>
			</div>
			
			<!-- Footer -->
			<div style="text-align: center; padding: 10px; background: #f1f1f1; font-size: 12px; color: #666;">
				본 메일은 자동 발송된 메일입니다. 문의사항이 있으시면 관리자에게 연락해주세요.
			</div>
		</div>
	</body>
	</html>
	`, req.Code)

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
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f9f9f9; padding: 20px; margin: 0;">
    <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border: 1px solid #ddd; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: #4caf50; color: white; padding: 20px; text-align: center;">
            <h1 style="margin: 0; font-size: 24px;">문의 등록 알림</h1>
        </div>
        
        <!-- Content -->
        <div style="padding: 20px;">
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">이름:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">이메일:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">휴대번호:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">내용:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
        </div>
        
        <!-- Footer -->
        <div style="text-align: center; padding: 10px; background: #f1f1f1; font-size: 12px; color: #666;">
            본 메일은 자동 발송된 메일입니다. 문의사항이 있으시면 관리자에게 연락해주세요.
        </div>
    </div>
</body>
</html>
`, req.Name, req.Email, req.Phone, req.Content)

	// 이메일 메시지 설정
	msg := []byte("To: kldga@haruharulab.com\r\n" +
		"Subject: kldga 문의등록 알림\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	// 이메일 전송
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, []string{"kldga@haruharulab.com"}, msg)
	if err != nil {
		return nil, err
	}

	return &pb.EmailResponse{Status: "Success"}, nil
}

func (s *EmailServer) KldgaSendCompetitionEmail(ctx context.Context, req *pb.KldgaCompetitionRequest) (*pb.EmailResponse, error) {
	// SMTP 설정
	email := os.Getenv("WELLKINSON_SMTP_EMAIL")
	password := os.Getenv("WELLKINSON_SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 인증 정보
	auth := smtp.PlainAuth("", email, password, smtpHost)
	gender := "남성"
	if req.Gender == 2 {
		gender = "여성"
	}
	// 이메일 본문 구성
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f9f9f9; padding: 20px; margin: 0;">
    <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border: 1px solid #ddd; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); overflow: hidden;">
        <!-- Header -->
        <div style="background: #4caf50; color: white; padding: 20px; text-align: center;">
            <h1 style="margin: 0; font-size: 24px;">대회 신청서 등록 알림</h1>
        </div>
        
        <!-- Content -->
        <div style="padding: 20px;">
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">이름:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">성별:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">신청리그:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
            <div style="margin-bottom: 15px;">
                <strong style="font-size: 16px; color: #4caf50;">현재 골프 전공 및 직업:</strong>
                <p style="margin: 0; font-size: 14px; color: #555;">%s</p>
            </div>
			<div style="margin-bottom: 15px;">
				<strong style="font-size: 16px; color: #4caf50;">휴대전화 번호:</strong>
				<p style="margin: 0; font-size: 14px; color: #555;">%s</p>
			</div>
			<div style="margin-bottom: 15px;">
				<strong style="font-size: 16px; color: #4caf50;">메모:</strong>
				<p style="margin: 0; font-size: 14px; color: #555;">%s</p>
			</div>
        </div>
        
        <!-- Footer -->
        <div style="text-align: center; padding: 10px; background: #f1f1f1; font-size: 12px; color: #666;">
            본 메일은 자동 발송된 메일입니다. 문의사항이 있으시면 관리자에게 연락해주세요.
        </div>
    </div>
</body>
</html>
`, req.Name, gender, req.League, req.Career, req.Phone, req.Memo)

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
