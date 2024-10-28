package core

import (
	"admin-service/model"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtSecretKey = []byte("adapfit_mark")

func verifyJWT(c *fiber.Ctx) (uint, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		log.Println("required")
		return 0, errors.New("authorization header is required")
	}

	// 'Bearer ' 접두사 제거
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("invalidtoken")
		return 0, errors.New("invalid token")
	}

	id := uint((*claims)["id"].(float64))

	return id, nil
}

func generateJWT(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), // 한달 유효 기간
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validatePhoneNumber(phone string) error {
	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil || !matched {
		log.Println("invalid format")
		return errors.New("invalid phone format, should be 01000000000")
	}
	return nil
}

func sendCode(number string) (string, error) {
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(10)) // 0부터 9까지의 숫자를 무작위로 선택
	}
	apiURL := "https://kakaoapi.aligo.in/akv10/alimtalk/send/"
	data := url.Values{}
	data.Set("apikey", os.Getenv("API_KEY"))
	data.Set("userid", os.Getenv("USER_ID"))
	data.Set("token", os.Getenv("TOKEN"))
	data.Set("senderkey", os.Getenv("SENDER_KEY"))
	data.Set("tpl_code", os.Getenv("TPL_CODE"))
	data.Set("sender", os.Getenv("SENDER"))
	data.Set("subject_1", os.Getenv("SUBJECT_1"))

	data.Set("receiver_1", number)
	data.Set("message_1", "인증번호는 ["+sb.String()+"]"+" 입니다.")

	// HTTP POST 요청 실행
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		fmt.Printf("HTTP Request Failed: %s\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println(fmt.Errorf("server returned non-200 status: %d, body: %s", resp.StatusCode, string(body)))

	return sb.String(), nil

}

func validateSignIn(request SignInRequest) error {
	// 빈 문자열 검사
	if request.Email == "" || len(request.Email) > 50 || strings.Contains(request.Email, " ") {
		log.Println("invalid1")
		return errors.New("invalid email format")
	}

	// 이메일 검증을 위한 정규 표현식
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(request.Email) {
		log.Println("invalid2")
		return errors.New("invalid email format")
	}

	name := strings.TrimSpace(request.Name)
	if utf8.RuneCountInString(name) > 5 || len(name) == 0 {
		log.Println("invalid3")
		return errors.New("invalid name")
	}

	major := strings.TrimSpace(request.Major)
	if len(major) == 0 {
		log.Println("invalid4")
		return errors.New("invalid major")
	}

	return nil
}
