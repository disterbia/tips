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

type UserRequest struct {
	ID           uint   `json:"-"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image" example:"base64string"`
	Phone        string `json:"phone"`
	Gender       bool   `json:"gender"`
	Birthday     string `json:"birthday" example:"yyyy-mm-dd"`
	UserType     uint   `json:"user_type"`
}

type UserResponse struct {
	Name         string           `json:"name"`
	Birthday     string           `json:"birthday" example:"yyyy-mm-dd"`
	Phone        string           `json:"phone"`
	Gender       bool             `json:"gender"` // true:남 false: 여
	SnsType      uint             `json:"sns_type"`
	CreatedAt    string           `json:"created_at"`
	ProfileImage ImageResponse    `json:"profile_image"`
	LinkedEmails []LinkedResponse `json:"linked_emails"`
}

type ImageResponse struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type LinkedResponse struct {
	SnsType uint   `json:"sns_type"`
	Email   string `json:"email"`
}

type LinkRequest struct {
	Id      uint   `json:"-"`
	IdToken string `json:"id_token"`
}

type AppVersionResponse struct {
	LatestVersion string `json:"latest_version"`
	AndroidLink   string `json:"android_link"`
	IosLink       string `json:"ios_link"`
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
