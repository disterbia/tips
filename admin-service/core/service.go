// /user-service/service/service.go

package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"admin-service/model"

	pb "admin-service/proto"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type AdminService interface {
	login(request LoginRequest) (string, error)
	searchHospitals(request SearchParam) ([]HospitalResponse, error)
	getPolicies() ([]PolicyResponse, error)
	verifyCode(target, code string) (string, error)
	sendAuthCodeForSignin(number string) (string, error)
	sendAuthCodeForId(r FindIdRequest) (string, error)
	signIn(r SignInRequest) (string, error)
	findId(r FindIdRequest) (string, error)
	sendAuthCodeForPw(r FindPwRequest) (string, error)
	changePw(r FindPasswordRequest) (string, error)
	question(r QuestionRequest) (string, error)
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
		log.Println("empty")
		return "", errors.New("empty")
	}

	// 이메일로 사용자 조회
	if err := service.db.Where("email = ? AND role_id = ?", request.Email, ADMINROLE).First(&u).Error; err != nil {
		log.Println("-2")
		return "", errors.New("-2")
	}

	// 비밀번호 비교
	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(request.Password)); err != nil {
		log.Println("-2")
		return "", errors.New("-2")
	}

	if !*u.IsApproval {
		log.Println("-1")
		return "", errors.New("-1")
	}

	// 새로운 JWT 토큰 생성
	tokenString, err := generateJWT(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *adminService) searchHospitals(request SearchParam) ([]HospitalResponse, error) {
	var hospitals []model.Hospital

	// 페이지네이션 설정
	pageSize := 50
	offset := int(request.Page) * pageSize

	keyword := request.Keyword

	// 키워드로 시작하는 결과를 우선 정렬하고, 그 뒤에 키워드가 포함된 결과를 정렬
	if err := service.db.Table("hospitals").
		Where("region_code = ? , name ILIKE ?", request.RegionCode, "%"+keyword+"%").
		Order("CASE WHEN name ILIKE '" + keyword + "%' THEN 1 ELSE 2 END").
		Limit(pageSize).
		Offset(offset).
		Find(&hospitals).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	var response []HospitalResponse
	for _, v := range hospitals {
		response = append(response, HospitalResponse{Name: v.Name, Number: v.Number})
	}

	return response, nil
}

func (service *adminService) getPolicies() ([]PolicyResponse, error) {
	var policies []model.Policy
	if err := service.db.Where("is_last = true").Find(&policies).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	var response []PolicyResponse
	for _, v := range policies {
		response = append(response, PolicyResponse{Title: v.Title, Body: v.Body})
	}

	return response, nil
}

func (service *adminService) sendAuthCodeForSignin(number string) (string, error) {

	err := validatePhoneNumber(number)
	if err != nil {
		return "", err
	}

	//존재하는 번호인지 체크
	result := service.db.Debug().Where("phone = ?", number).First(&model.User{})
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("db error")
			return "", errors.New("db error")
		}

	} else if result.RowsAffected > 0 {
		// 레코드가 존재할 때
		log.Println("-1")
		return "", errors.New("-1")
	}

	code, err := sendCode(number)

	if err != nil {
		return "", err
	}

	if err := service.db.Create(&model.AuthCode{Phone: number, Code: code}).Error; err != nil {
		return "", err
	}
	return "200", nil
}

func (service *adminService) verifyCode(target, code string) (string, error) {
	if target == "" || code == "" {
		log.Println("-1")
		return "", errors.New("-1")
	}

	now := time.Now()
	threeMinutesAgo := now.Add(-3 * time.Minute)
	var authCode model.AuthCode

	if err := service.db.Where("(phone = ? OR email = ?) AND created_at >= ? ", target, target, threeMinutesAgo).Last(&authCode).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	if authCode.Code != code {
		log.Println("-1")
		return "", errors.New("-1")
	}

	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if err := tx.Where("(phone = ? OR email = ?)", target, target).Unscoped().Delete(&model.AuthCode{}).Error; err != nil {
		tx.Rollback()
		log.Println("db error3")
		return "", errors.New("db error3")
	}

	if err := tx.Create(&model.VerifiedTarget{Target: target}).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	tx.Commit()

	return "200", nil
}

func (service *adminService) signIn(r SignInRequest) (string, error) {
	if err := validateSignIn(r); err != nil {
		log.Println("-21")
		return "", errors.New("-2")
	}

	birthday, err := time.Parse("2006-01-02", r.Birthday)
	if err != nil {
		log.Println("-22")
		return "", errors.New("-2")
	}
	// 비밀번호 공백 제거
	password := strings.TrimSpace(r.Password)

	if password == "" {
		log.Println("-23")
		return "", errors.New("-2")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	now := time.Now()
	thirtyMinutesAgo := now.Add(-30 * time.Minute)

	if err := service.db.Where("target = ? AND created_at >= ?", r.Phone, thirtyMinutesAgo).Last(&model.VerifiedTarget{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-1")
			return "", errors.New("-1") // 인증해야함
		}
		log.Println("db error")
		return "", errors.New("db error")
	}

	finalPassword := string(hashedPassword)
	falseValue := false
	var user = model.User{Name: r.Name, Email: &r.Email, Password: &finalPassword, Phone: r.Phone, Birthday: birthday, HospitalID: &r.HospitalID, Major: &r.Major,
		RoleID: uint(ADMINROLE), IsApproval: &falseValue}

	if err := service.db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			log.Println("-3")
			return "", errors.New("-3")
		}
		log.Println("db error2")
		return "", errors.New("db error2")
	}
	if err := service.db.Unscoped().Where("target = ? ", r.Phone).Delete(&model.VerifiedTarget{}); err != nil {
		log.Println("db error9")
		return "", errors.New("db error9")
	}
	return "200", nil
}

func (service *adminService) sendAuthCodeForId(r FindIdRequest) (string, error) {

	birthday, err := time.Parse("2006-01-02", r.Birthday)
	if err != nil {
		log.Println("-1")
		return "", errors.New("-1")
	}

	var user model.User
	if err := service.db.Where("name = ? AND phone = ? AND birthday = ? AND role_id = ?", r.Name, r.Phone, birthday, ADMINROLE).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-1")
			return "", errors.New("-1") // 찾을수 없음
		}
		log.Println("db error")
		return "", errors.New("db error")
	}

	code, err := sendCode(r.Phone)

	if err != nil {
		return "", err
	}

	if err := service.db.Create(&model.AuthCode{Phone: r.Phone, Code: code}).Error; err != nil {
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	return "200", nil
}

func (service *adminService) findId(r FindIdRequest) (string, error) {

	birthday, err := time.Parse("2006-01-02", r.Birthday)
	if err != nil {
		log.Println("-2")
		return "", errors.New("-2")
	}

	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	if err := service.db.Where("target = ? AND created_at >= ?", r.Phone, oneMinuteAgo).Last(&model.VerifiedTarget{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-1")
			return "", errors.New("-1") // 인증해야함
		}
		log.Println("db error")
		return "", errors.New("db error")
	}

	var user model.User
	if err := service.db.Where("name = ? AND phone = ? AND birthday = ? AND role_id = ?", r.Name, r.Phone, birthday, ADMINROLE).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-2")
			return "", errors.New("-2") // 찾을수 없음
		}
		log.Println("db error")
		return "", errors.New("db error2")
	}

	if err := service.db.Unscoped().Where("target = ? ", r.Phone).Delete(&model.VerifiedTarget{}); err != nil {
		log.Println("db error9")
		return "", errors.New("db error9")
	}
	return *user.Email, nil
}

func (service *adminService) sendAuthCodeForPw(r FindPwRequest) (string, error) {

	var sb strings.Builder
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10)) // 0부터 9까지의 숫자를 무작위로 선택
	}

	if r.Phone == "" {
		if err := service.db.Where("email = ? AND role_id = ? ", r.Email, ADMINROLE).First(&model.User{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("-1")
				return "", errors.New("-1") // 찾을수 없음
			}
			log.Println("db error")
			return "", errors.New("db error")
		}
		response, err := service.emailClient.SendCodeEmail(context.Background(), &pb.EmailCodeRequest{
			Email: r.Email, // 받는 사람의 이메일
			Code:  sb.String(),
		})

		if err != nil {
			log.Println(err)
			return "", err
		}

		if response != nil && response.Status == "Success" {
			if err := service.db.Create(&model.AuthCode{Email: r.Email, Code: sb.String()}).Error; err != nil {
				log.Println(err)
				return "", errors.New("db error2")
			}
		}

		log.Printf("send email: %v", response)

	} else {

		if err := service.db.Where("email = ? AND phone = ? AND role_id = ? ", r.Email, r.Phone, ADMINROLE).First(&model.User{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("-1")
				return "", errors.New("-1") // 찾을수 없음
			}
			log.Println("db error3")
			return "", errors.New("db error3")
		}
		//phone 인증번호 전송
		code, err := sendCode(r.Phone)

		if err != nil {
			log.Println(err)
			return "", err
		}

		if err := service.db.Create(&model.AuthCode{Phone: r.Phone, Code: code}).Error; err != nil {
			log.Println("db error4")
			return "", errors.New("db error4")
		}
	}

	return "200", nil
}

func (service *adminService) changePw(r FindPasswordRequest) (string, error) {
	// 비밀번호 공백 제거
	password := strings.TrimSpace(r.Password)

	if password == "" {
		log.Println("-2")
		return "", errors.New("-2")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	now := time.Now()
	threeMinutesAgo := now.Add(-3 * time.Minute)

	if err := service.db.Where("( target = ? OR target = ? ) AND created_at >= ?", r.Phone, r.Email, threeMinutesAgo).Last(&model.VerifiedTarget{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-1")
			return "", errors.New("-1") // 인증해야함
		}
		log.Println("db error")
		return "", errors.New("db error")
	}

	finalPassword := string(hashedPassword)

	if err := service.db.Model(&model.User{}).Where("email = ?", r.Email).UpdateColumn("password", finalPassword).Error; err != nil {
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	if err := service.db.Unscoped().Where("target = ? OR target = ? ", r.Phone, r.Email).Delete(&model.VerifiedTarget{}); err != nil {
		log.Println("db error9")
		return "", errors.New("db error9")
	}
	return "200", nil
}

func (service *adminService) question(r QuestionRequest) (string, error) {
	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, r.Phone)
	if err != nil || !matched {
		log.Println("invalid phone")
		return "", errors.New("invalid phone format, should be 01000000000")
	}

	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(r); err != nil {
		log.Println(err)
		return "", err
	}

	question := model.Question{Email: r.Email, Phone: r.Phone, Name: r.Name, PossibleTime: r.PossibleTime, EntryRoute: r.EntryRoute, HospitalName: r.HospitalName}
	if err := service.db.Create(&question).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}

	return "200", nil
}
