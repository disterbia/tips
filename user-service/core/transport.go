package core

import (
	"context"
	"sync"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
)

var userLocks sync.Map

// @Tags 로그인 /user
// @Summary sns 로그인
// @Description sns 로그인 성공시 호출
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "요청 DTO - idToken 필수, user- user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자 / 최초 로그인 이후 로그인시 fcm_token,device_id 만 필요함"
// @Success 200 {object} SuccessResponse "성공시 JWT 토큰 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 인증필요 , "-2" = 이미 가입한 번호,  "-3" = 추가정보 입력 필요"
// @Router /sns-login [post]
func SnsLoginHandler(loginEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := loginEndpoint(context.Background(), req)
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
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 인증필요 , "-2" = 이미 가입한 번호,  "-3" = 추가정보 입력 필요"
// @Router /phone-login [post]
func PhoneLoginHandler(loginEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PhoneLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		response, err := loginEndpoint(context.Background(), req)
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
// @Description 인증번호 발송시 호출
// @Accept  json
// @Produce  json
// @Param number path string true "휴대번호"
// @Success 200 {object} BasicResponse "성공시 200 반환"
// @Failure 400 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} ErrorResponse "요청 처리 실패시 오류 메시지 반환: 오류메시지 "-1" = 이미 가입한번호"
// @Router /send-code/{number} [post]
func SendCodeHandler(sendEndpoint endpoint.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		number := c.Params("email")

		response, err := sendEndpoint(c.Context(), number)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

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
// @Param request body UserRequest true "요청 DTO - 업데이트 할 데이터/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자"
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
// @Success 200 {object} UserResponse "성공시 유저 객체 반환/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자 sns_type- 0:휴대폰,1:카카오 2:구글 3:애플"
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
// @Success 200 {object} dto.BasicResponse "성공시 200 반환"
// @Failure 400 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
// @Failure 500 {object} dto.ErrorResponse "요청 처리 실패시 오류 메시지 반환"
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
