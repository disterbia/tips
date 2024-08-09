package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"medicine-service/model"
	"sort"
	"time"

	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type MedicineService interface {
	SaveMedicine(r MedicineRequest) (string, error)
	RemoveMedicines(id uint, uid uint) (string, error)
	GetExpects(id uint) ([]MedicineTakeResponse, error)
	GetMedicines(id uint) ([]MedicineResponse, error)
	TakeMedicine(takeMedicine TakeMedicine) (string, error)
	UnTakeMedicine(id, uid uint) (string, error)
	SearchMedicines(keyword string) ([]string, error)
}

type medicineService struct {
	db   *gorm.DB
	nats *nats.Conn
}

func NewMedicineService(db *gorm.DB, nats *nats.Conn) MedicineService {
	return &medicineService{db: db, nats: nats}

}

func (service *medicineService) SaveMedicine(r MedicineRequest) (string, error) {

	var weekdays []uint
	var times []string
	var weekdaysResult pq.Int64Array
	var startAt *time.Time
	var endAt *time.Time
	var minReserves *float32
	var remaining *float32

	for _, v := range r.Weekdays {
		if v == 0 || v > 7 {
			return "", errors.New("validate weekdays")
		}
		weekdays = append(weekdays, v)
	}
	if weekdays != nil {
		weekdaysResult = uintSliceToInt64Array(weekdays)
	}

	for _, v := range r.Times {
		time, err := time.Parse("15:04", v)
		if err != nil {
			return "", errors.New("validate times")
		}
		times = append(times, time.Format("15:04"))
	}

	if r.StartAt != "" {
		time, err := time.Parse("2006-01-02", r.StartAt)
		if err != nil {
			return "", errors.New("invalid date format, should be YYYY-MM-DD")
		}
		startAt = &time
	}
	if r.EndAt != "" {
		time, err := time.Parse("2006-01-02", r.EndAt)
		if err != nil {
			return "", errors.New("invalid date format, should be YYYY-MM-DD")
		}
		endAt = &time
	}

	if r.MinReserves != 0 {
		minReserves = &r.MinReserves
	}
	if r.Remaining != 0 {
		remaining = &r.Remaining
	}

	var medicine model.Medicine
	if err := service.db.Where("id = ? AND uid = ? ", r.Id, r.Uid).First(&medicine).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("db error")
		}
	}

	medicine.Uid = r.Uid
	medicine.Name = r.Name
	medicine.Weekdays = weekdaysResult
	medicine.Times = times
	medicine.Dose = r.Dose
	medicine.MedicineType = r.MedicineType
	medicine.UsePrivacy = r.UsePrivacy
	medicine.StartAt = startAt
	medicine.EndAt = endAt
	medicine.MinReserves = minReserves
	medicine.Remaining = remaining
	medicine.IsActive = r.IsActive

	if err := service.db.Save(&medicine).Error; err != nil {
		return "", errors.New("db error2")
	}

	body := medicine.Name + " " + fmt.Sprintf("%v", medicine.Dose) + " " + medicine.MedicineType + " 먹을 시간입니다. 드시고 나면 잊지 말고 표시해주세요."
	if medicine.UsePrivacy {
		body = "약 먹을 시간입니다. 드시고 나면 잊지 말고 표시해주세요."
	}

	if medicine.IsActive {
		notificationRequest := NotificationRequest{
			Uid:      r.Uid,
			Type:     uint(MEDICINENOTIFICATION),
			ParentID: medicine.ID,
			Body:     body,
			StartAt:  r.StartAt,
			EndAt:    r.EndAt,
			Times:    times,
			Week:     r.Weekdays,
		}
		eventData, err := json.Marshal(notificationRequest)
		if err != nil {
			return "", errors.New("json marshal error")
		}

		if err := service.nats.Publish("save-notification", eventData); err != nil {
			return "", errors.New("nats publish error")
		}
	} else {
		notificationRequest := NotificationRequest{
			Uid:      r.Uid,
			Type:     uint(MEDICINENOTIFICATION),
			ParentID: medicine.ID,
		}
		eventData, err := json.Marshal(notificationRequest)
		if err != nil {
			return "", errors.New("json marshal error")
		}

		if err := service.nats.Publish("remove-notification", eventData); err != nil {
			return "", errors.New("nats publish error")
		}
	}

	return "200", nil
}

func (service *medicineService) RemoveMedicines(id uint, uid uint) (string, error) {
	if err := service.db.Where("id =? ", id).Delete(&model.Medicine{}).Error; err != nil {
		return "", errors.New("db error")
	}

	// nats removeNotifications
	return "200", nil
}

func (service *medicineService) GetExpects(uid uint) ([]MedicineTakeResponse, error) {
	var responses []MedicineTakeResponse
	var user model.User

	// 1. 사용자의 정보를 데이터베이스에서 가져옵니다.
	if err := service.db.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 2. 사용자의 의약품 리스트를 데이터베이스에서 가져옵니다.
	var medicines []model.Medicine
	if err := service.db.Where("uid = ? AND is_active = ?", uid, true).Find(&medicines).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 3. 사용자의 의약품 복용 기록을 데이터베이스에서 가져옵니다.
	var medicineTakes []model.MedicineTake
	if err := service.db.Where("uid = ?", uid).Find(&medicineTakes).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 4. 복용 기록을 날짜와 시간으로 맵핑합니다.
	//    예: medicineTakeMap["2024-08-08"]["08:00"] = 1 (복용 기록 ID)
	medicineTakeMap := make(map[string]map[string]uint)
	for _, take := range medicineTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		if _, exists := medicineTakeMap[dateStr]; !exists {
			medicineTakeMap[dateStr] = make(map[string]uint)
		}
		medicineTakeMap[dateStr][take.TimeTaken] = take.ID
	}

	// 5. 사용자의 가입 날짜부터 오늘까지의 날짜 범위를 생성합니다.
	startDate := user.CreatedAt
	endDate := time.Now()
	// 날짜별 응답을 저장할 맵
	responseMap := make(map[string]*MedicineTakeResponse)
	// 6. 각 의약품에 대해 복용 스케줄을 생성합니다.
	for _, medicine := range medicines {
		medStart := medicine.StartAt
		medEnd := medicine.EndAt

		// 6.1. 의약품 시작 날짜가 사용자의 가입 날짜 이전이면 가입 날짜로 설정합니다.
		if medStart == nil || medStart.Before(startDate) {
			medStart = &startDate
		}

		// 6.2. 의약품 종료 날짜가 오늘 날짜 이후이면 오늘 날짜로 설정합니다.
		if medEnd == nil || medEnd.After(endDate) {
			medEnd = &endDate
		}

		// 7. 의약품의 복용 기간 동안 반복합니다.
		for d := *medStart; !d.After(*medEnd); d = d.AddDate(0, 0, 1) {
			weekDay := int(d.Weekday())

			// 7.1. 해당 날짜의 요일이 의약품의 복용 요일에 포함되는지 확인합니다.
			if contains(medicine.Weekdays, int64(weekDay)) {
				dateStr := d.Format("2006-01-02")
				timeTaken := make(map[string]*uint)

				// 7.2. 의약품의 각 복용 시간에 대해 복용 기록을 확인합니다.
				for _, timeStr := range medicine.Times {
					var takeId *uint
					if val, exists := medicineTakeMap[dateStr][timeStr]; exists {
						takeId = &val
					} else {
						takeId = nil
					}
					timeTaken[timeStr] = takeId
				}

				// 7.3. 해당 날짜에 대한 ExpectMedicineResponse를 생성합니다.
				response := ExpectMedicineResponse{
					Id:        medicine.ID,
					Name:      medicine.Name,
					TimeTaken: timeTaken,
					Dose:      medicine.Dose,
				}

				// 7. 해당 날짜의 응답을 맵에서 가져오거나 새로 만듭니다.
				if _, exists := responseMap[dateStr]; !exists {
					responseMap[dateStr] = &MedicineTakeResponse{
						DateTaken:     dateStr,
						MedicineTaken: []ExpectMedicineResponse{},
					}
				}

				// 8. 해당 날짜의 응답에 추가합니다.
				responseMap[dateStr].MedicineTaken = append(responseMap[dateStr].MedicineTaken, response)
			}
		}
	}

	// 맵을 리스트로 변환
	for _, response := range responseMap {
		responses = append(responses, *response)
	}
	// 8. 모든 복용 기록을 처리하여 응답에 포함합니다.
	for _, take := range medicineTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		timeStr := take.TimeTaken

		var found bool

		// 8.1. 이미 응답에 해당 날짜의 기록이 있는지 확인합니다.
		for i, res := range responses { // 1. responses 리스트에서 각 응답(res)을 순회합니다.
			if res.DateTaken == dateStr { // 2. 응답의 날짜가 현재 복용 기록의 날짜와 같은지 확인합니다.
				for j, medRes := range res.MedicineTaken { // 3. 해당 날짜의 MedicineTaken 리스트를 순회합니다.
					if medRes.Id == take.MedicineID { // 4. 해당 MedicineTaken의 ID가 현재 복용 기록의 MedicineID와 같은지 확인합니다.
						responses[i].MedicineTaken[j].TimeTaken[timeStr] = &take.ID // 5. 같은 약물 기록이 있으면, 해당 시간의 복용 기록 ID를 업데이트합니다.
						found = true                                                // 6. 업데이트가 완료되었음을 표시합니다.
						break                                                       // 7. 내부 루프를 탈출합니다.
					}
				}
			}
			if found {
				break // 8. 외부 루프를 탈출합니다.
			}
		}

		// 8.2. 해당 날짜와 시간에 대한 기록이 없으면 새로 추가합니다.
		if !found {
			timeTaken := map[string]*uint{
				timeStr: &take.ID,
			}
			response := ExpectMedicineResponse{
				Id:        take.MedicineID,
				Name:      take.Medicine.Name,
				TimeTaken: timeTaken,
				Dose:      take.Dose,
			}
			dateResponse := MedicineTakeResponse{
				DateTaken:     dateStr,
				MedicineTaken: []ExpectMedicineResponse{response},
			}
			responses = append(responses, dateResponse)
		}
	}

	// responses를 DateTaken 순으로 정렬
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].DateTaken < responses[j].DateTaken
	})

	return responses, nil
}

func (service *medicineService) GetMedicines(id uint) ([]MedicineResponse, error) {
	var medicines []model.Medicine

	if err := service.db.Where("uid = ?", id).Find(&medicines).Error; err != nil {
		return nil, errors.New("db error")
	}

	var medicineResponses []MedicineResponse

	for _, v := range medicines {
		var startAt, endAt *string
		var minReserves, remaining *float32

		if v.StartAt != nil {
			start := v.StartAt.Format("2006-01-02")
			startAt = &start
		}
		if v.EndAt != nil {
			end := v.EndAt.Format("2006-01-02")
			endAt = &end
		}
		if v.Remaining != nil {
			remaining = v.Remaining
		}
		if v.MinReserves != nil {
			minReserves = v.MinReserves
		}
		log.Println(v.ID)
		log.Println(v.CreatedAt)
		log.Println(v.Name)
		medicineResponses = append(medicineResponses, MedicineResponse{Id: v.ID, Name: v.Name, Times: v.Times, Weekdays: int64ArrayToUintSlice(v.Weekdays), Dose: v.Dose, MedicineType: v.MedicineType,
			StartAt: startAt, EndAt: endAt, MinReserves: minReserves, Remaining: remaining, UsePrivacy: v.UsePrivacy, IsActive: v.IsActive})
	}

	return medicineResponses, nil
}

func (service *medicineService) TakeMedicine(request TakeMedicine) (string, error) {
	var medicineTake model.MedicineTake
	var medicine model.Medicine

	_, err := time.Parse("15:04", request.TimeTaken)
	if err != nil {
		return "", errors.New("invalid time format, should be HH:MM")
	}

	dateTaken, err := time.Parse("2006-01-02", request.DateTaken)
	if err != nil {
		return "", errors.New("invalid date format, should be YYYY-MM-DD")
	}

	result := service.db.Where("medicine_id = ? AND uid = ? AND date_taken = ? AND time_taken = ?",
		request.MedicineId, request.Uid, dateTaken, request.TimeTaken).First(&medicineTake)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("db error")
		}
	}
	if result.RowsAffected != 0 {
		return "", errors.New("already taken")
	}

	if err := service.db.Where("id = ?", request.MedicineId).First(&medicine).Error; err != nil {
		return "", errors.New("db error2")
	}

	if medicine.Remaining != nil {
		if request.Dose > *medicine.Remaining {
			return "", errors.New("over remaining")
		}
	}

	tx := service.db.Begin()
	if medicine.Remaining != nil {
		if err := tx.UpdateColumn("remaining", *medicine.Remaining-request.Dose).Error; err != nil {
			tx.Rollback()
			return "", errors.New("db error3")
		}
	}

	medicineTake = model.MedicineTake{MedicineID: request.MedicineId, Uid: request.Uid, Dose: request.Dose, DateTaken: dateTaken, TimeTaken: request.TimeTaken}

	if err := tx.Create(&medicineTake).Error; err != nil {
		tx.Rollback()
		return "", errors.New("db error4")
	}

	tx.Commit()
	return "200", nil
}

func (service *medicineService) UnTakeMedicine(id, uid uint) (string, error) {

	var medicineTake model.MedicineTake
	var medicine model.Medicine

	if err := service.db.Where("id = ? AND uid = ?", id, uid).First(&medicineTake).Error; err != nil {
		return "", errors.New("db error")
	}
	tx := service.db.Begin()

	result := tx.Where("id = ?").Unscoped().Delete(&model.MedicineTake{})
	if result.Error != nil {
		tx.Rollback()
		return "", errors.New("db error")
	}

	if result.RowsAffected != 0 {
		if err := tx.Model(&medicine).Where("id = ?", medicineTake.MedicineID).UpdateColumn("remaining", gorm.Expr("remaining + ?", medicineTake.Dose)).Error; err != nil {
			tx.Rollback()
			return "", errors.New("db error2")
		}
	}

	tx.Commit()

	return "200", nil
}

func (service *medicineService) SearchMedicines(keyword string) ([]string, error) {
	var names []string
	err := service.db.Model(&model.MedicineSearch{}).Where("name LIKE ?", "%"+keyword+"%").Pluck("name", &names).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	return names, nil
}
