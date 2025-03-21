package core

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

var ipLimiters = make(map[string]*rate.Limiter)
var ipLimitersMutex sync.Mutex

func getClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return c.IP()
}

// @Tags 로그인 /user
// @Summary sns 로그인
// @Description sns 로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "요청 DTO - idToken 필수, user- user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자 / 최초 로그인 이후 로그인시 fcm_token,device_id 만 필요함"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 PHONE=0, KAKAO=1, GOOGLE=2, APPLE=3 / -1 = 인증필요 , -2 = 추가정보 입력 필요 "
// @Router /sns-login [post]
func SnsLoginHandler(loginEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := loginEndpoint(c.Context(), req)
		resp := response.(LoginResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"jwt": resp.Jwt})
	}
}

// @Tags 로그인 /user
// @Summary 휴대번호 로그인
// @Description 휴대번호 로그인시 호출
// @Accept  json
// @Produce  json
// @Param request body PhoneLoginRequest true "요청 DTO user- user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자 / 최초 로그인 이후 로그인시 phone,fcm_token,device_id 만 필요함"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 KAKAO=1, GOOGLE=2, APPLE=3 / -1 = 인증필요 , -2 = 추가정보 입력 필요 "
// @Router /phone-login [post]
func PhoneLoginHandler(loginEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PhoneLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := loginEndpoint(c.Context(), req)
		resp := response.(LoginResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"jwt": resp.Jwt})
	}
}

// @Tags 로그인 /user
// @Summary 자동로그인
// @Description 최초 로그인 이후 앱 실행시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body AutoLoginRequest true "요청 DTO"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Security jwt
// @Router /auto-login [post]
func AutoLoginHandler(autoLoginEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var req AutoLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		req.Id = id
		response, err := autoLoginEndpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(LoginResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 인증번호 /user
// @Summary 인증번호 발송
// @Description 회원가입 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param number path string true "휴대번호"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 이미 가입한번호"
// @Router /send-code-join/{number} [post]
func SendCodeForSignInHandler(sendEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		number := c.Params("number")

		// IP 주소를 가져오기 위한 함수 호출
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(24*time.Hour), 10)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청이 허용되지 않으면 에러 반환
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}
		response, err := sendEndpoint(c.Context(), number)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 응답이 성공적이면 RateLimiter를 업데이트
		ipLimitersMutex.Lock()
		limiter.Allow()
		ipLimitersMutex.Unlock()

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 인증번호 /user
// @Summary 인증번호 발송
// @Description 휴대번호 로그인 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param number path string true "휴대번호"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /send-code-login/{number} [post]
func SendCodeForLoginHandler(sendEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		number := c.Params("number")

		// IP 주소를 가져오기 위한 함수 호출
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(24*time.Hour), 10)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청이 허용되지 않으면 에러 반환
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}
		response, err := sendEndpoint(c.Context(), number)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 응답이 성공적이면 RateLimiter를 업데이트
		ipLimitersMutex.Lock()
		limiter.Allow()
		ipLimitersMutex.Unlock()

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 인증번호 /user
// @Summary 번호 인증
// @Description 인증번호 입력 후 호출
// @Accept  json
// @Produce  json
// @Param request body VerifyRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 코드불일치"
// @Router /verify-code [post]
func VerifyHandler(verifyEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req VerifyRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := verifyEndpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 마이페이지 /user
// @Summary 내정보 변경
// @Description 내정보 변경시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body UserRequest true "요청 DTO - 업데이트 할 데이터/ ture:남성 user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자 "
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 번호인증 필요"
// @Router /update-user [post]
func UpdateUserHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var req UserRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		req.ID = id

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}
		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 회원정보 조회 /user
// @Summary 회원정보 조회
// @Description 내 정보 조회시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} UserResponse "성공시 유저 객체 반환/ ture:남성 user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자  sns_type- 0:휴대폰,1:카카오 2:구글 3:애플"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-user [get]
func GetUserHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}
		resp := response.(UserResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 계정 연동 /user
// @Summary 계정 연동
// @Description 계정 연동시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body LinkRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /link-email [post]
func LinkHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var req LinkRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		req.Id = id
		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}
		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 회원탈퇴 /user
// @Summary 회원탈퇴
// @Description 회원탈퇴시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-user [post]
func RemoveHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 공통 /user
// @Summary 최신버전 조회
// @Description 최신버전 조회시 호출
// @Accept  json
// @Produce  json
// @Success 200 {object} AppVersionResponse "최신 버전 정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-version [get]
func GetVersionHandeler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		response, err := endpoint(c.Context(), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(AppVersionResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 공통 /user
// @Summary 약관 조회
// @Description 약관 조회시 호출
// @Accept  json
// @Produce  json
// @Success 200 {object} []PoliceResponse "약관 정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-polices [get]
func GetPolicesHandeler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		response, err := endpoint(c.Context(), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}
		resp := response.([]PoliceResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

func AppleCallbackHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("POST!!!")
		code := c.Query("code")
		state := c.Query("state")
		log.Println("code:", code, "state:", state)
		// POST 요청에서 body 파싱
		var req AppleCallbackRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if req.Code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "authorization code is missing",
			})
		}

		// 엔드포인트 호출
		response, err := endpoint(context.Background(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// 응답에서 ID 토큰 추출
		resp := response.(AppleCallbackResponse)
		idToken := resp.IDToken
		log.Println("ID Token:", idToken)

		// ID 토큰을 앱 딥 링크로 리다이렉트
		appScheme := "myapp://callback" // 앱의 딥 링크 URI 스킴
		redirectURL := fmt.Sprintf("%s?id_token=%s&code=%s", appScheme, idToken, code)

		// 앱으로 리다이렉트 (302 리다이렉트)
		return c.Redirect(redirectURL, fiber.StatusFound)
	}
}

// func AppleCallbackHandler(endpoint endpoint.Endpoint) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		log.Println("POST!!!")
// 		code := c.Query("code")
// 		state := c.Query("state")
// 		log.Println("code:", code, "state:", state)
// 		// POST 요청에서 body 파싱
// 		var req AppleCallbackRequest
// 		if err := c.BodyParser(&req); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": "invalid request body",
// 			})
// 		}

// 		if req.Code == "" {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": "authorization code is missing",
// 			})
// 		}

// 		// 엔드포인트 호출
// 		response, err := endpoint(context.Background(), req)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}

// 		resp := response.(AppleCallbackResponse)
// 		log.Println("idToken:", resp.IDToken)
// 		return c.JSON(fiber.Map{
// 			"access_token": resp.AccessToken,
// 			"id_token":     resp.IDToken,
// 		})
// 	}
// }
