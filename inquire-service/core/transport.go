package core

import (
	"strconv"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 로그인 /inquire
// @Summary 관리자 로그인
// @Description 관리자 로그인시 호출
// @Accept  json
// @Produce  json
// @Param email body string true "email"
// @Param password body string true "password"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /login [post]
func AdminLoginHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req map[string]interface{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		email := req["email"].(string)
		password := req["password"].(string)
		response, err := endpoint(c.Context(), map[string]interface{}{
			"email":    email,
			"password": password,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

		}

		resp := response.(LoginResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 문의 /inquire
// @Summary 답변/추가문의
// @Description 답변/추가문의 등록시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body InquireReplyRequest true "요청 DTO - 답변데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /inquire-reply [post]
func AnswerHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)
		var req InquireReplyRequest
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

// @Tags 문의 /inquire
// @Summary 문의하기
// @Description 문의등록시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body InquireRequest true "요청 DTO - 문의데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /send-inquire [post]
func SendHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)

		var req InquireRequest
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

// @Tags 문의 /inquire
// @Summary 문의조회(본인)
// @Description 나의문의보기시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  page  query  uint  false  "페이지 번호 default 0" (30개씩)
// @Param  start_date  query string  false  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  false  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []InquireResponse "문의내역 배열 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-inquires [get]
func GetHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var queryParams GetInquireParams
		if err := c.QueryParser(&queryParams); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// id와 queryParams를 함께 전달
		response, err := endpoint(c.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]InquireResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 문의 /inquire
// @Summary 문의조회(관리자)
// @Description 관리자 문의내역 확인시 호출 (30개씩)
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param  page  query  uint  false  "페이지 번호 default 0"
// @Param  start_date  query string  false  "시작날짜 yyyy-mm-dd"
// @Param  end_date  query string  false  "종료날짜 yyyy-mm-dd"
// @Success 200 {object} []InquireResponse "문의내역 배열 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /all-inquires [get]
func GetAllHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var queryParams GetInquireParams
		if err := c.QueryParser(&queryParams); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// id와 queryParams를 함께 전달
		response, err := endpoint(c.Context(), map[string]interface{}{
			"id":          id,
			"queryParams": queryParams,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]InquireResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 문의 /inquire
// @Summary 문의삭제
// @Description 문의삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path string ture "문의ID"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-inquire/{id} [post]
func RemoveInquireHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		alarmId := c.Params("id")
		id, err := strconv.Atoi(alarmId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		response, err := endpoint(c.Context(), map[string]interface{}{
			"uid": uid,
			"id":  uint(id),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 문의 /inquire
// @Summary 문의답변/추가문의 삭제
// @Description 문의답변/추가문의 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path string ture "답변/추가문의ID"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-reply/{id} [post]
func RemoveReplyHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		alarmId := c.Params("id")
		id, err := strconv.Atoi(alarmId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		response, err := endpoint(c.Context(), map[string]interface{}{
			"uid": uid,
			"id":  uint(id),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
