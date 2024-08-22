package core

type GetEmotionsParams struct {
	StartDate string `query:"start_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
}
type EmotionRequest struct {
	Uid        uint   `json:"-"`
	Emotion    uint   `json:"emotion" validate:"required,min=1,max=5"`
	Symptoms   []uint `json:"symptoms"`
	Memo       string `json:"memo" validate:"max=500"`
	TargetDate string `json:"target_date" example:"YYYY-mm-dd" validate:"required,datetime=2006-01-02"`
}

type EmotionResponse struct {
	Emotion    uint   `json:"emotion"`
	Symptoms   []uint `json:"symptoms"`
	Memo       string `json:"memo"`
	TargetDate string `json:"target_date" example:"YYYY-mm-dd"`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}
type ErrorResponse struct {
	Err string `json:"err"`
}

type BasicResponse struct {
	Code string `json:"code"`
}
