package core

import (
	"strconv"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 운동 /exercise
// @Summary 운동 저장
// @Description 운동 등록 및 수정시 호출 - 생성시 id생략
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body ExerciseRequest true "요청 DTO - 운동 데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /save-exercise [post]
func SaveExerciseHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)

		var req ExerciseRequest
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

// @Tags 운동 /exercise
// @Summary 운동 내역 조회
// @Description 운동 내역 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []ExerciseTakeResponse "운동 내역정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-takens [get]
func GetExpectsHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.([]ExerciseTakeResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 운동 /exercise
// @Summary 등록 운동 조회
// @Description 등록 운동 조회시 호출
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Success 200 {object} []ExerciseResponse "등록 운동 정보"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /get-exercises [get]
func GetExercisesHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := endpoint(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		resp := response.([]ExerciseResponse)
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}

// @Tags 운동 /exercise
// @Summary 운동 삭제
// @Description 운동 삭제시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body []uint true "삭제할 id 배열"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /remove-exercise/{id} [post]
func RemoveExercisesHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		exerId := c.Params("id")
		eid, err := strconv.Atoi(exerId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		response, err := endpoint(c.Context(), map[string]interface{}{
			"uid": id,
			"id":  uint(eid),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)

	}
}

// @Tags 운동 /exercise
// @Summary 운동 기록
// @Description 운동 완료 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param request body TakeExercise true "운동 완료 데이터"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /do-exercises [post]
func DoExerciseHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)

		var req TakeExercise
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

// @Tags 운동 /exercise
// @Summary 운동 기록
// @Description 운동 취소시 호출
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {jwt_token}"
// @Param id path string ture "취소 ID"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Router /undo-exercise/{id} [post]
func UnDoHandler(endpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := verifyJWT(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// 사용자별 잠금 시작
		if _, loaded := userLocks.LoadOrStore(id, true); loaded {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concurrent request detected"})
		}
		defer userLocks.Delete(id)

		doId := c.Params("id")
		did, err := strconv.Atoi(doId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		response, err := endpoint(c.Context(), map[string]interface{}{
			"uid": id,
			"id":  uint(did),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp := response.(BasicResponse)
		return c.Status(fiber.StatusOK).JSON(resp)

	}
}

// // @Tags 운동 /exercise
// // @Summary 운동 동영상 카테고리 조회
// // @Description 운동 동영상 카테고리 조회시 호출
// // @Produce  json
// // @Param Authorization header string true "Bearer {jwt_token}"
// // @Success 200 {object} []ProjectResponse "카테고리 정보"
// // @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// // @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// // @Router /get-projects [get]
// func GetProjectsHandler(endpoint endpoint.Endpoint) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		_, err := verifyJWT(c)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		response, err := endpoint(c.Context(), nil)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		resp := response.([]ProjectResponse)
// 		return c.Status(fiber.StatusOK).JSON(resp)
// 	}
// }

// // @Tags 운동 /exercise
// // @Summary 카테고리별 운동 동영상 조회 (20개씩)
// // @Description 카테고리별 운동 동영상 조회시 호출
// // @Produce  json
// // @Param Authorization header string true "Bearer {jwt_token}"
// // @Param  project_id  query string  true  "project_id"
// // @Param  page  query uint  false  "페이지 default 0"
// // @Success 200 {object} []VideoResponse "동영상 정보"
// // @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// // @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// // @Router /get-videos [get]
// func GetVideosHandler(endpoint endpoint.Endpoint) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		_, err := verifyJWT(c)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		var queryParams GetVideoParams
// 		if err := c.QueryParser(&queryParams); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		response, err := endpoint(c.Context(), queryParams)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		resp := response.([]VideoResponse)
// 		return c.Status(fiber.StatusOK).JSON(resp)
// 	}
// }
