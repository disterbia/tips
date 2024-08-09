package core

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type GetInquireParams struct {
	Page      uint   `query:"page"`
	StartDate string `query:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `query:"end_date" example:"YYYY-MM-DD"`
}

type InquireRequest struct {
	Uid     uint   `json:"-"`
	Email   string `json:"email" validate:"required,email"`
	Title   string `json:"title"  validate:"required,max=50"`
	Content string `json:"content"  validate:"required,max=1000"`
}

type InquireResponse struct {
	Id        uint                   `json:"id"`
	Email     string                 `json:"email"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	CreatedAt string                 `json:"created_at" example:"YYYY-mm-dd HH:mm:ss"`
	Replies   []InquireReplyResponse `json:"replies"`
}

type InquireReplyRequest struct {
	Uid       uint   `json:"-"`
	InquireId uint   `json:"inquire_id"`
	Content   string `json:"content" validate:"required,max=1000"`
	ReplyType bool   `json:"reply_type"`
}

type InquireReplyResponse struct {
	Id        uint   `json:"id"`
	InquireId uint   `json:"inquire_id"`
	Content   string `json:"content"`
	ReplyType bool   `json:"reply_type"`
	Created   string `json:"created" example:"YYYY-mm-dd HH:mm:ss "`
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
