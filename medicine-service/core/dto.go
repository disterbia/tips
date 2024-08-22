package core

type GetParams struct {
	StartDate string `query:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `query:"end_date" example:"YYYY-MM-DD"`
}

type MedicineRequest struct {
	Id           uint     `json:"id"`
	Uid          uint     `json:"-"`
	Name         string   `json:"name"`
	Weekdays     []uint   `json:"weekdays"`
	Times        []string `json:"times" example:"HH:mm,HH:mm"`
	Dose         float32  `json:"dose"`
	MedicineType string   `json:"medicine_type"`
	StartAt      string   `json:"start_at" example:"YYYY-MM-dd"`
	EndAt        string   `json:"end_at"  example:"YYYY-MM:dd"`
	Remaining    float32  `json:"remaining"`
	UsePrivacy   bool     `json:"use_privacy"`
	IsActive     bool     `json:"is_active"`
}
type MedicineResponse struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	Times        []string `json:"times"`
	Weekdays     []uint   `json:"weekdays"`
	Dose         float32  `json:"dose"`
	MedicineType string   `json:"medicine_type"`
	StartAt      *string  `json:"start_at" example:"YYYY-MM-dd"`
	EndAt        *string  `json:"end_at"  example:"YYYY-MM:dd"`
	Remaining    float32  `json:"remaining"`
	UsePrivacy   bool     `json:"use_privacy"`
	IsActive     bool     `json:"is_active"`
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

type MedicineTakeResponse struct {
	DateTaken     string                   `json:"date_taken" example:"YYYY-MM-dd"`
	MedicineTaken []ExpectMedicineResponse `json:"medicine_taken"`
}

type ExpectMedicineResponse struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	TimeTaken map[string]*uint `json:"time_taken"`
	Dose      float32          `json:"dose"`
}

type TakeMedicine struct {
	Uid        uint    `json:"-"`
	MedicineId uint    `json:"medicine_id"`
	DateTaken  string  `json:"date_taken"  example:"YYYY-MM-DD"`
	TimeTaken  string  `json:"time_taken"  example:"HH:mm"`
	Dose       float32 `json:"dose"`
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
