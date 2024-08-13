package core

import (
	"encoding/json"
	"errors"
	"exercise-service/model"
	"sort"
	"time"

	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ExerciseService interface {
	SaveExercise(r ExerciseRequest) (string, error)
	GetExpects(id uint) ([]ExerciseTakeResponse, error)
	RemoveExercise(id uint, uid uint) (string, error)
	DoExercise(exerciseDo TakeExercise) (string, error)
	UndoExercise(id, uid uint) (string, error)
	GetExercises(id uint) ([]ExerciseResponse, error)
	GetProjects() ([]ProjectResponse, error)
	GetVideos(projectId string, page uint) ([]VideoResponse, error)
}

type exerciseService struct {
	db   *gorm.DB
	nats *nats.Conn
}

func NewExerciseService(db *gorm.DB, nats *nats.Conn) ExerciseService {
	return &exerciseService{db: db, nats: nats}
}

func (service *exerciseService) SaveExercise(r ExerciseRequest) (string, error) {

	var weekdays []uint
	var times []string
	var weekdaysResult pq.Int64Array
	var startAt *time.Time
	var endAt *time.Time

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

	var exercise model.Exercise
	if err := service.db.Where("id = ? AND uid = ? ", r.Id, r.Uid).First(&exercise).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("db error")
		}
	}

	exercise.Uid = r.Uid
	exercise.Name = r.Name
	exercise.Weekdays = weekdaysResult
	exercise.Times = times
	exercise.StartAt = startAt
	exercise.EndAt = endAt
	exercise.IsActive = r.IsActive

	if err := service.db.Save(&exercise).Error; err != nil {
		return "", errors.New("db error2")
	}

	body := "운동 할 시간입니다!"

	if exercise.IsActive {
		notificationRequest := NotificationRequest{
			Uid:      r.Uid,
			Type:     uint(EXERCISENOTIFIATION),
			ParentID: exercise.ID,
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

		if err := service.nats.Publish("save-exercise", eventData); err != nil {
			return "", errors.New("nats publish error")
		}
	} else {
		notificationRequest := NotificationRequest{
			Uid:      r.Uid,
			Type:     uint(EXERCISENOTIFIATION),
			ParentID: exercise.ID,
		}
		eventData, err := json.Marshal(notificationRequest)
		if err != nil {
			return "", errors.New("json marshal error")
		}

		if err := service.nats.Publish("remove-exercise", eventData); err != nil {
			return "", errors.New("nats publish error")
		}
	}

	return "200", nil
}

func (service *exerciseService) RemoveExercise(id uint, uid uint) (string, error) {
	if err := service.db.Where("id =? ", id).Delete(&model.Exercise{}).Error; err != nil {
		return "", errors.New("db error")
	}

	// nats removeNotifications
	return "200", nil
}

func (service *exerciseService) GetExpects(uid uint) ([]ExerciseTakeResponse, error) {
	var responses []ExerciseTakeResponse
	var user model.User

	// 1. 사용자의 정보를 데이터베이스에서 가져옵니다.
	if err := service.db.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 2. 사용자의 의약품 리스트를 데이터베이스에서 가져옵니다.
	var exercises []model.Exercise
	if err := service.db.Where("uid = ? AND is_active = ?", uid, true).Find(&exercises).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 3. 사용자의 의약품 복용 기록을 데이터베이스에서 가져옵니다.
	var exerciseTakes []model.ExerciseTake
	if err := service.db.Where("uid = ?", uid).Find(&exerciseTakes).Error; err != nil {
		return nil, errors.New("db error")
	}

	// 4. 복용 기록을 날짜와 시간으로 맵핑합니다.
	//    예: exerciseTakeMap["2024-08-08"]["08:00"] = 1 (복용 기록 ID)
	exerciseTakeMap := make(map[string]map[string]uint)
	for _, take := range exerciseTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		if _, exists := exerciseTakeMap[dateStr]; !exists {
			exerciseTakeMap[dateStr] = make(map[string]uint)
		}
		exerciseTakeMap[dateStr][take.TimeTaken] = take.ID
	}

	// 5. 사용자의 가입 날짜부터 오늘까지의 날짜 범위를 생성합니다.
	startDate := user.CreatedAt
	endDate := time.Now()
	// 날짜별 응답을 저장할 맵
	responseMap := make(map[string]*ExerciseTakeResponse)
	// 6. 각 의약품에 대해 복용 스케줄을 생성합니다.
	for _, exercise := range exercises {
		exerStart := exercise.StartAt
		exerEnd := exercise.EndAt

		// 6.1. 의약품 시작 날짜가 사용자의 가입 날짜 이전이면 가입 날짜로 설정합니다.
		if exerStart == nil || exerStart.Before(startDate) {
			exerStart = &startDate
		}

		// 6.2. 의약품 종료 날짜가 오늘 날짜 이후이면 오늘 날짜로 설정합니다.
		if exerEnd == nil || exerEnd.After(endDate) {
			exerEnd = &endDate
		}

		// 7. 의약품의 복용 기간 동안 반복합니다.
		for d := *exerStart; !d.After(*exerEnd); d = d.AddDate(0, 0, 1) {
			weekDay := int(d.Weekday())

			// 7.1. 해당 날짜의 요일이 의약품의 복용 요일에 포함되는지 확인합니다.
			if contains(exercise.Weekdays, int64(weekDay)) {
				dateStr := d.Format("2006-01-02")
				timeTaken := make(map[string]*uint)

				// 7.2. 의약품의 각 복용 시간에 대해 복용 기록을 확인합니다.
				for _, timeStr := range exercise.Times {
					var takeId *uint
					if val, exists := exerciseTakeMap[dateStr][timeStr]; exists {
						takeId = &val
					} else {
						takeId = nil
					}
					timeTaken[timeStr] = takeId
				}

				// 7.3. 해당 날짜에 대한 ExpectMedicineResponse를 생성합니다.
				response := ExpectExerciseResponse{
					Id:        exercise.ID,
					Name:      exercise.Name,
					TimeTaken: timeTaken,
				}

				// 7. 해당 날짜의 응답을 맵에서 가져오거나 새로 만듭니다.
				if _, exists := responseMap[dateStr]; !exists {
					responseMap[dateStr] = &ExerciseTakeResponse{
						DateTaken:     dateStr,
						ExerciseTaken: []ExpectExerciseResponse{},
					}
				}

				// 8. 해당 날짜의 응답에 추가합니다.
				responseMap[dateStr].ExerciseTaken = append(responseMap[dateStr].ExerciseTaken, response)
			}
		}
	}

	// 맵을 리스트로 변환
	for _, response := range responseMap {
		responses = append(responses, *response)
	}
	// 8. 모든 복용 기록을 처리하여 응답에 포함합니다.
	for _, take := range exerciseTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		timeStr := take.TimeTaken

		var found bool

		// 8.1. 이미 응답에 해당 날짜의 기록이 있는지 확인합니다.
		for i, res := range responses { // 1. responses 리스트에서 각 응답(res)을 순회합니다.
			if res.DateTaken == dateStr { // 2. 응답의 날짜가 현재 복용 기록의 날짜와 같은지 확인합니다.
				for j, exerRes := range res.ExerciseTaken { // 3. 해당 날짜의 MedicineTaken 리스트를 순회합니다.
					if exerRes.Id == take.ExerciseID { // 4. 해당 MedicineTaken의 ID가 현재 복용 기록의 MedicineID와 같은지 확인합니다.
						responses[i].ExerciseTaken[j].TimeTaken[timeStr] = &take.ID // 5. 같은 약물 기록이 있으면, 해당 시간의 복용 기록 ID를 업데이트합니다.
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
			response := ExpectExerciseResponse{
				Id:        take.ExerciseID,
				Name:      take.Exercise.Name,
				TimeTaken: timeTaken,
			}
			dateResponse := ExerciseTakeResponse{
				DateTaken:     dateStr,
				ExerciseTaken: []ExpectExerciseResponse{response},
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

func (service *exerciseService) GetExercises(id uint) ([]ExerciseResponse, error) {
	var exercises []model.Exercise

	if err := service.db.Where("uid = ?", id).Find(&exercises).Error; err != nil {
		return nil, errors.New("db error")
	}

	var exerciseResponses []ExerciseResponse

	for _, v := range exercises {
		var startAt, endAt *string

		if v.StartAt != nil {
			start := v.StartAt.Format("2006-01-02")
			startAt = &start
		}
		if v.EndAt != nil {
			end := v.EndAt.Format("2006-01-02")
			endAt = &end
		}

		exerciseResponses = append(exerciseResponses, ExerciseResponse{Id: v.ID, Name: v.Name, Times: v.Times, Weekdays: int64ArrayToUintSlice(v.Weekdays),
			StartAt: startAt, EndAt: endAt, IsActive: v.IsActive})
	}

	return exerciseResponses, nil
}

func (service *exerciseService) DoExercise(request TakeExercise) (string, error) {

	_, err := time.Parse("15:04", request.TimeTaken)
	if err != nil {
		return "", errors.New("invalid time format, should be HH:MM")
	}

	dateTaken, err := time.Parse("2006-01-02", request.DateTaken)
	if err != nil {
		return "", errors.New("invalid date format, should be YYYY-MM-DD")
	}

	var exerciseTake model.ExerciseTake

	result := service.db.Where("exercise_id = ? AND uid = ? AND date_taken = ? AND time_taken = ?",
		request.ExerciseId, request.Uid, dateTaken, request.TimeTaken).First(&exerciseTake)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("db error")
		}
	}
	if result.RowsAffected != 0 {
		return "", errors.New("already taken")
	}

	exerciseTake = model.ExerciseTake{ExerciseID: request.ExerciseId, Uid: request.Uid, DateTaken: dateTaken, TimeTaken: request.TimeTaken}

	if err := service.db.Create(&exerciseTake).Error; err != nil {
		return "", errors.New("db error2")
	}

	return "200", nil
}

func (service *exerciseService) UndoExercise(id, uid uint) (string, error) {

	if err := service.db.Where("id = ? AND uid = ?", id, uid).Unscoped().Delete(&model.ExerciseTake{}).Error; err != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}

func (service *exerciseService) GetProjects() ([]ProjectResponse, error) {

	var projects []ProjectResponse
	err := service.db.Model(&model.Video{}).
		Select("project_id, project_name as name, count(*) as count").
		Group("project_id").Scan(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

// face-service 의 face-exercise 쪽에는 한번에 다가져옴

func (service *exerciseService) GetVideos(projectId string, page uint) ([]VideoResponse, error) {
	pageSize := 20
	var videos []model.Video
	offset := page * uint(pageSize)
	if err := service.db.Where("project_id = ?", projectId).Offset(int(offset)).Limit(pageSize).Order("id DESC").Find(&videos).Error; err != nil {
		return nil, errors.New("db error")
	}

	var videoResponses []VideoResponse
	for _, v := range videos {
		videoResponses = append(videoResponses, VideoResponse{Name: v.Name, VideoId: v.VideoId, ThumbnailUrl: v.ThumbnailUrl, Duration: v.Duration})
	}

	return videoResponses, nil
}
