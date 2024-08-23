package core

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SearchParam struct {
	Keyword    string `json:"keyword"`
	RegionCode string `json:"region_code" query:"region_code"`
	Page       uint   `json:"page"`
}

type HospitalResponse struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type PolicyResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type SignInRequest struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday" example:"yyyy-mm-dd"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	HospitalID uint   `json:"hospital_id"`
	Major      string `json:"major"`
}

type FindIdRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday" example:"yyyy-mm-dd"`
}

type FindPwRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone" example:"이메일로 찾기시 생략"`
}

type FindPasswordRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone" example:"이메일로 찾기시 생략"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number" example:"01000000000"`
	Code        string `json:"code" example:"인증번호 6자리"`
}

type QuestionRequest struct {
	Name         string `json:"name" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	HospitalName string `json:"hospital_name" validate:"required"`
	PossibleTime string `json:"possible_time" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	EntryRoute   string `json:"entry_route"  validate:"required"`
}

type BasicResponse struct {
	Code string `json:"code"`
}

// // for swagger ////
type SuccessResponse struct {
	Jwt string `json:"jwt"`
}
type ErrorResponse struct {
	Err string `json:"err"`
}
