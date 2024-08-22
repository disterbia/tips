package core

import (
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

// @Tags 회원가입 /admin
// @Summary 회원가입
// @Description 관리자 회원가입시 호출
// @Accept  json
// @Produce  json
// @Param request body SignInRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 -1: 인증필요 , -2: invalid body , -3: 이미 가입된 이메일"
// @Router /sign-in [post]
func SignInHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SignInRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 로그인 /admin
// @Summary 로그인
// @Description 로그인시 호출
// @Accept  json
// @Produce  json
// @Param email body string true "email"
// @Param password body string true "password"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 -1: 승인필요 , -2: 이메일/비밀번호 틀림"
// @Router /login [post]
func LoginHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(LoginResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 회원가입 /admin
// @Summary 병원검색
// @Description 병원검색시 호출
// @Accept  json
// @Produce  json
// @Param name body string true "name"
// @Param page body int true "page default 0"
// @Param region_code body string true "region_code"
// @Success 200 {object} []HospitalResponse "병원정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /search-hospitals [get]
func SearchHospitalsHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var param SearchParam

		if err := c.QueryParser(&param); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		log.Println(param.Keyword)
		log.Println(param.Page)
		log.Println(param.RegionCode)

		response, err := endpoint(c.Context(), param)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.([]HospitalResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 회원가입 /admin
// @Summary 이용약관 가져오기
// @Description 이용약관 내용 조회시 호출
// @Accept  json
// @Produce  json
// @Success 200 {object} []PolicyResponse "정책정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-policies [get]
func GetPoliciesHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		response, err := endpoint(c.Context(), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]PolicyResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 회원가입 /admin
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

// @Tags 아이디 찾기 /admin
// @Summary 인증번호 발송
// @Description 아이디 찾기 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param request body FindIdRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 오류메시지 "-1" 정보 불일치"
// @Router /send-code-id [post]
func SendCodeForIdHandler(sendEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FindIdRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

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
		response, err := sendEndpoint(c.Context(), req)
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

// @Tags 인증번호 인증 /admin
// @Summary 인증번호 인증
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

// @Tags 비밀번호 찾기 /admin
// @Summary 인증번호 발송
// @Description 비밀번호 찾기 시 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param request body FindPwRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 오류메시지 "-1" 정보 불일치"
// @Router /send-code-pw [post]
func SendCodeForPwHandler(sendEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FindPwRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
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
		response, err := sendEndpoint(c.Context(), req)
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

// @Tags 아이디 찾기 /admin
// @Summary 아이디 찾기
// @Description 아이디 찾기시 호출
// @Accept  json
// @Produce  json
// @Param request body FindIdRequest true "요청 DTO"
// @Success 200 {object} string "성공시 email 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 -1: 인증필요 , -2: invalid pw"
// @Router /find-id [post]
func FindIdHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FindIdRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(string)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 비밀번호 찾기 /admin
// @Summary 비밀번호 변경
// @Description 비밀번호 변경시 호출
// @Accept  json
// @Produce  json
// @Param request body FindPasswordRequest true "요청 DTO"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환 -1: 인증필요 , -2: invalid pw"
// @Router /change-pw [post]
func ChangePwHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FindPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
