package core

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

// @Tags 알람 /notification
// @Summary 수신 알림 조회
// @Description 수신 알림 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []MessageResponse "알람정보 - type 1:수면 2:운동"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-messages [get]
func GetMessagesHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]MessageResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 알람 /notification
// @Summary  알림 모두 읽기
// @Description 수신 알림 조회시 자동 읽음 처리라면 수신 알림 조회완료 후 함께 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /read-messages [post]
func ReadAllHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

// @Tags 알람 /notification
// @Summary 수신된 알림 삭제
// @Description 수신된 알림 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-messages [post]
func RemoveMessagesHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var ids []uint // 삭제할 ID 배열
		if err := c.BodyParser(&ids); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		response, err := endpoint(c.Context(), map[string]interface{}{
			"uid": uid,
			"ids": ids,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)

	}
}
