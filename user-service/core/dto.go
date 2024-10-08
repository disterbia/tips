package core

type LoginRequest struct {
	IdToken  string `json:"id_token"`
	DeviceID string `json:"device_id"`
	FCMToken string `json:"fcm_token"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday" example:"yyyy-mm-dd"`
	Gender   bool   `json:"gender"`
	UserType uint   `json:"user_type"`
}

type PhoneLoginRequest struct {
	Phone    string `json:"phone"`
	DeviceID string `json:"device_id"`
	FCMToken string `json:"fcm_token"`
	Name     string `json:"name"`
	Birthday string `json:"birthday" example:"yyyy-mm-dd"`
	Gender   bool   `json:"gender"`
	UserType uint   `json:"user_type"`
}

type AutoLoginRequest struct {
	Id       uint   `json:"-"`
	FcmToken string `json:"fcm_token"`
	DeviceId string `json:"device_id"`
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
	UserType     uint             `json:"user_type"`
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

type PoliceResponse struct {
	PoliceType uint   `json:"police_type"`
	Title      string `json:"title"`
	Body       string `json:"body"`
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
