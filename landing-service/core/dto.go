package core

type KldgaInquireRequest struct {
	Name    string `json:"name" validate:"required,max=50"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"  validate:"required,max=11"`
	Content string `json:"content"  validate:"required,max=1000"`
}

type KldgaCompetitionRequest struct {
	Name   string `json:"name" validate:"required,max=50"`
	Gender uint   `json:"gender" validate:"required,min=1,max=2"`
	League string `json:"league" validate:"required,max=50"`
	Career string `json:"career" validate:"required,max=50"`
	Phone  string `json:"phone"  validate:"required,max=11"`
	Memo   string `json:"memo"  validate:"max=100"`
}

type AdapfitInquireReqeust struct {
	Name    string `json:"name" validate:"required,max=50"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"  validate:"required,max=11"`
	Purpose string `json:"purpose"  validate:"required,max=1000"`
	Career  string `json:"career" validate:"required,max=50"`
	Content string `json:"content"  validate:"required,max=1000"`
}

type AuthCodeRequest struct {
	Phone string `json:"phone"`
}

type VerifyAuthRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

type ErrorResponse struct {
	Err string `json:"err"` // wwwwww
}

type BasicResponse struct {
	Code string `json:"code"`
}
