package core

import (
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

var ipLimiters = make(map[string]*rate.Limiter)
var ipLimiters2 = make(map[string]*rate.Limiter)
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

// @Tags 랜딩페이지 /landing
// @Summary kldga 문의하기
// @Description 문의 등록시 호출
// @Accept  json
// @Produce  json
// @Param request body KldgaInquireRequest true "요청 DTO - 문의데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /kldga/inquire [post]
func KldgaInquireHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req KldgaInquireRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// IP 주소를 가져오기 위한 함수 호출
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(24*time.Hour), 5)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청이 허용되지 않으면 에러 반환
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}
		response, err := endpoint(c.Context(), req)
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

// @Tags 랜딩페이지 /landing
// @Summary kldga 대회 신청서 작성
// @Description kldga 대회 신청서 등록시 호출
// @Accept  json
// @Produce  json
// @Param request body KldgaCompetitionRequest true "요청 DTO - 신청서 데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /kldga/inquire [post]
func KldgaCompetitionHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req KldgaCompetitionRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// IP 주소를 가져오기 위한 함수 호출
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(24*time.Hour), 5)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청이 허용되지 않으면 에러 반환
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}
		response, err := endpoint(c.Context(), req)
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

// @Tags 랜딩페이지 /landing
// @Summary kldga 인증번호 전송
// @Description 대회 신청서 작성 인증번호 전송시 호출
// @Accept  json
// @Produce  json
// @Param request body AuthCodeRequest true "요청 DTO - 휴대번호"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /kldga/send-code [post]
func SendAuthCodeHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req AuthCodeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// IP 주소를 가져오기 위한 함수 호출
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters2[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(24*time.Hour), 5)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		// 요청이 허용되지 않으면 에러 반환
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}
		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 응답이 성공적이면 RateLimiter를 업데이트
		ipLimitersMutex.Lock()
		limiter.Allow()
		ipLimitersMutex.Unlock()
		return c.Status(fiber.StatusOK).JSON(response)
	}

}

// @Tags 랜딩페이지 /landing
// @Summary kldga 인증번호 인증
// @Description 대회 신청서 작성 인증번호 인증시 호출
// @Accept  json
// @Produce  json
// @Param request body VerifyAuthRequest true "요청 DTO - 인증정보"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /kldga/verify-code [post]
func VerifyAuthCodeHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req VerifyAuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
