package core

type GetFaceInfoParams struct {
	StartDate string `query:"start_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
}

type GetTapBlinkScoreParams struct {
	ScoreType uint   `query:"score_type" validate:"required,min=1,max=2"`
	StartDate string `query:"start_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" example:"yyyy-mm-dd" validate:"required,datetime=2006-01-02"`
}

type SampleVideoResponse struct {
	Category  uint   `json:"category"`
	VideoType uint   `json:"video_type"`
	Title     string `json:"title"`
	VideoId   string `json:"video_id"`
}

type FaceInfoResponse struct {
	TargetDate string              `json:"date" example:"YYYY-mm-dd"`
	FaceInfos  map[uint][]FaceInfo `json:"face_infos"`
}
type FaceInfo struct {
	FaceLine uint    `json:"face_line"`
	Sd       float64 `json:"sd"`
}

type TapBlinkResponse struct {
	TargetDate    string  `json:"date" example:"YYYY-mm-dd"`
	ScoreType     uint    `json:"score_type"`
	SuccessCount  uint    `json:"success_count"`
	ErrorCount    uint    `json:"error_count"`
	ReactionSpeed float64 `json:"reaction_speed"`
}

type FaceInfoRequest struct {
	FaceInfos map[uint][]FaceInfo `json:"face_infos"`
}

type TapBlinkRequest struct {
	Uid           uint    `json:"-"`
	ScoreType     uint    `json:"score_type" validate:"required,min=1,max=2"`
	SuccessCount  uint    `json:"success_count" validate:"required,max=100"`
	ErrorCount    uint    `json:"error_count" validate:"required,max=100"`
	ReactionSpeed float64 `json:"reaction_speed" validate:"required,gt=0,max=100"`
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
