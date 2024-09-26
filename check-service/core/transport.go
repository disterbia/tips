package core

import (
	"log"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 검사 /check
// @Summary 샘플동영상 전체 조회
// @Description 샘플동영상 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []SampleVideoResponse "동영상 정보 - category: 1-표정 2-손가락태핑 3-눈깜빡임 / video_type: 1-태핑,눈깜빡임,기쁨 2-슬픔 3-놀람 4-분노"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-videos [get]
func GetSampleVideosHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		_, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// id와 queryParams를 함께 전달
		response, err := endpoint(c.Context(), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]SampleVideoResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 검사 /check
// @Summary 표정검사 점수 조회
// @Description 표정검사 점수 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  start_date  query string  true  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  true  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []FaceScoreResponse "점수정보 - face_type: 1-기쁨 2-슬픔 3-놀람 4-분노"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-face-scores [get]
func GetFaceScoresHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var queryParams []GetFaceScoreParams

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

		resp := response.([]FaceScoreResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 검사 /check
// @Summary 손가락태핑/눈깜빡임 점수 조회
// @Description 손가락태핑/눈깜빡임 점수 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  score_type query uint  true  "1-손가락태핑 2-눈깜빡임"
// @Param  start_date query string  true  "시작날짜 yyyy-mm-dd"
// @Param  end_date query string  true  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []TapBlinkResponse "점수정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-scores [get]
func GetScoresHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 토큰 검증 및 처리
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var queryParams []GetTapBlinkScoreParams

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

		resp := response.([]TapBlinkResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 검사 /check
// @Summary 표정검사 점수 저장
// @Description 표정검사 완료 후 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body FaceScoreRequest true "요청 DTO - face_type: 1-기쁨 2-슬픔 3-놀람 4-분노"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-face-score [post]
func SaveFaceScoreHandler(endpoint endpoint.Endpoint) fiber.Handler {
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
		var req FaceScoreRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": req,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 검사 /check
// @Summary 손가락태핑/눈깜빡임 검사 점수 저장
// @Description 손가락태핑/눈깜빡임 검사 완료 후 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body TapBlinkRequest true "요청 DTO - score_type: 1-손가락태핑 2-눈깜빡임"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-score [post]
func SaveScoreHandler(endpoint endpoint.Endpoint) fiber.Handler {
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
		var req TapBlinkRequest
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
