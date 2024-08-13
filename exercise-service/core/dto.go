package core

type GetParams struct {
	StartDate string `form:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `form:"end_date" example:"YYYY-MM-DD"`
}

type GetVideoParams struct {
	Page      uint   `form:"page"`
	ProjectId string `form:"project_id"`
}

type ExerciseRequest struct {
	Id       uint     `json:"id"`
	Uid      uint     `json:"-"`
	Name     string   `json:"name"`
	StartAt  string   `json:"start_at" example:"YYYY-MM-dd"`
	EndAt    string   `json:"end_at"  example:"YYYY-MM:dd"`
	Times    []string `json:"times" example:"HH:mm,HH:mm"`
	Weekdays []uint   `json:"weekdays"`
	IsActive bool     `json:"is_active"`
}

type ExerciseResponse struct {
	Id       uint     `json:"id"`
	Name     string   `json:"name"`
	Times    []string `json:"times"`
	Weekdays []uint   `json:"weekdays"`
	StartAt  *string  `json:"start_at" example:"YYYY-MM-dd"`
	EndAt    *string  `json:"end_at"  example:"YYYY-MM:dd"`
	IsActive bool     `json:"is_active"`
}

type ExerciseTakeResponse struct {
	DateTaken     string                   `json:"date_taken" example:"YYYY-MM-dd"`
	ExerciseTaken []ExpectExerciseResponse `json:"Exercise_taken"`
}

type ExpectExerciseResponse struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	TimeTaken map[string]*uint `json:"time_taken"`
}

type TakeExercise struct {
	Uid        uint   `json:"-"`
	ExerciseId uint   `json:"exercise_id"`
	DateTaken  string `json:"date_taken"  example:"YYYY-MM-DD"`
	TimeTaken  string `json:"time_taken"  example:"HH:mm"`
}

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

type ProjectResponse struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
	Count     uint   `json:"count"`
}

type VideoResponse struct {
	Name         string `json:"name"`
	VideoId      string `json:"video_id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Duration     uint   `json:"duration"`
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
