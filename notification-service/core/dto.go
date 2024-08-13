package core

type NotificationRequest struct {
	Uid      uint
	Type     uint
	ParentID uint
	Body     string
	StartAt  string
	EndAt    string
	Times    []string
	Week     []uint
}
type MessageResponse struct {
	Id        uint   `json:"id"`
	Type      uint   `json:"type"`
	Body      string `json:"body"`
	ParentId  uint   `json:"parent_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" example:"YYYY-mm-dd HH:mm:ss "`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

type ErrorResponse struct {
	Err string `json:"err" `
}

type BasicResponse struct {
	Code string `json:"code"`
}
