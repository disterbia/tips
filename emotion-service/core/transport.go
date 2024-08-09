package core

import (
	"log"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 기분 /emotion
// @Summary 기분 생성/수정
// @Description 기분 생성시 Id 생략
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body EmotionRequest true "요청 DTO - 기분 데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-emotion [post]
func SaveEmotionHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)
		var req EmotionRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		req.Uid = id
		response, err := endpoint(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 기분 /emotion
// @Summary 기분 조회
// @Description 기분 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  true  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  true  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []EmotionResponse "기분정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-emotions [get]
func GetEmotionsHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var queryParams GetEmotionsParams

		if err := c.QueryParser(&queryParams); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		log.Println(queryParams)
		// id와 queryParams를 함께 전달
		response, err := endpoint(c.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]EmotionResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
