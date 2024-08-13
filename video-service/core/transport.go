package core

import (
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 관리자 동영상 관리 /video
// @Summary 최상위 레벨 조회
// @Description 최초에 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []VimeoLevel1 "웰킨스 폴더 내용"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-items [get]
func GetVimeoLevel1sHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		resp := response.([]VimeoLevel1)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 관리자 동영상 관리 /video
// @Summary 폴더 레벨2 조회
// @Description 폴더내부 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path string true "id"
// @Success 200 {object} []VimeoLevel2 "해당 폴더 내용"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-items/{id} [get]
func GetVimeoLevel2sHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		projectId := c.Query("id")
		response, err := endpoint(c.Context(), map[string]interface{}{
			"id":        id,
			"projectId": projectId,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]VimeoLevel2)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 관리자 동영상 관리 /video
// @Summary 동영상 활성화
// @Description 활성화 동영상 변경시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body VideoData true "활성화 할 id 배열"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-videos/{id} [post]
func SaveHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var videoData VideoData // 삭제할 ID 배열
		if err := c.BodyParser(&videoData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)

		videoData.Id = id
		response, err := endpoint(c.Context(), videoData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
