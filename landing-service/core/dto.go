package core

type KldgaRequest struct {
	Name    string `json:"name" validate:"required,max=50"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"  validate:"required,max=11"`
	Content string `json:"content"  validate:"required,max=1000"`
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
