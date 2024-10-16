package core

import (
	"encoding/json"
	"errors"
	"exercise-service/model"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ExerciseService interface {
	saveExercise(r ExerciseRequest) (string, error)
	getExpects(id uint) ([]ExerciseTakeResponse, error)
	removeExercise(id uint, uid uint) (string, error)
	doExercise(exerciseDo TakeExercise) (string, error)
	undoExercise(id, uid uint) (string, error)
	getExercises(id uint) ([]ExerciseResponse, error)
	// getProjects() ([]ProjectResponse, error)
	// getVideos(projectId string, page uint) ([]VideoResponse, error)
}

type exerciseService struct {
	db   *gorm.DB
	nats *nats.Conn
}

func NewExerciseService(db *gorm.DB, nats *nats.Conn) ExerciseService {
	return &exerciseService{db: db, nats: nats}
}

func (service *exerciseService) saveExercise(r ExerciseRequest) (string, error) {

	var weekdays []uint
	var times []string
	var weekdaysResult pq.Int64Array
	var startAt *time.Time
	var endAt *time.Time

	name := strings.TrimSpace(r.Name)
	if utf8.RuneCountInString(name) > 10 || len(name) == 0 {
		return "", errors.New("validate name")
	}

	for _, v := range r.Weekdays {
		if v > 6 {
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

func (service *exerciseService) removeExercise(id uint, uid uint) (string, error) {
	if err := service.db.Where("id =? AND uid =?", id, uid).Delete(&model.Exercise{}).Error; err != nil {
		return "", errors.New("db error")
	}

	notificationRequest := NotificationRequest{
		Uid:      uid,
		Type:     uint(EXERCISENOTIFIATION),
		ParentID: id,
	}

	eventData, err := json.Marshal(notificationRequest)
	if err != nil {
		return "", errors.New("json marshal error")
	}

	if err := service.nats.Publish("remove-exercise", eventData); err != nil {
		return "", errors.New("nats publish error")
	}

	return "200", nil
}

func (service *exerciseService) getExpects(uid uint) ([]ExerciseTakeResponse, error) {
	var responses []ExerciseTakeResponse
	var user model.User

	if err := service.db.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, errors.New("db error")
	}

	var exercises []model.Exercise
	if err := service.db.Where("uid = ? ", uid).Find(&exercises).Error; err != nil {
		return nil, errors.New("db error")
	}

	var exerciseTakes []model.ExerciseTake
	// exercise 만 unscoped 먹이기 exerciseTakes는 제외
	if err := service.db.Where("uid = ?", uid).Preload("Exercise", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Find(&exerciseTakes).Error; err != nil {
		return nil, errors.New("db error")
	}

	//  예: exerciseTakeMap["2024-08-08"]["08:00"] = 1 (복용 기록 ID)
	exerciseTakeMap := make(map[uint]map[string]map[string]uint)
	for _, take := range exerciseTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		if _, exists := exerciseTakeMap[take.ExerciseID]; !exists {
			exerciseTakeMap[take.ExerciseID] = make(map[string]map[string]uint)
		}
		if _, exists := exerciseTakeMap[take.ExerciseID][dateStr]; !exists {
			exerciseTakeMap[take.ExerciseID][dateStr] = make(map[string]uint)
		}
		exerciseTakeMap[take.ExerciseID][dateStr][take.TimeTaken] = take.ID
	}

	// 사용자의 가입 날짜부터 오늘까지의 날짜 범위를 생성합니다.
	startDate := user.CreatedAt.Truncate(24 * time.Hour)
	endDate := time.Now().Truncate(24 * time.Hour)
	// 날짜별 응답을 저장할 맵
	responseMap := make(map[string]ExerciseTakeResponse)

	for _, exercise := range exercises {
		exerStart := exercise.StartAt
		exerEnd := exercise.EndAt

		if exerStart == nil || exerStart.Before(startDate) {
			exerStart = &startDate
		}
		if exerEnd == nil || exerEnd.After(endDate) {
			exerEnd = &endDate
		}
		// 운동이 활성화되지 않았고, 해당 날짜 이후에는 보여주지 않도록 처리

		for d := *exerStart; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
			weekDay := int(d.Weekday())
			dateStr := d.Format("2006-01-02")
			// 비활성화된 운동은 수행 기록이 없는 경우 날짜에 표시되지 않도록 처리
			if !exercise.IsActive {
				if _, exists := exerciseTakeMap[exercise.ID][dateStr]; !exists {
					continue
				}
			}
			if contains(exercise.Weekdays, int64(weekDay)) {
				timeTaken := make(map[string]*uint)

				for _, timeStr := range exercise.Times {
					if val, exists := exerciseTakeMap[exercise.ID][dateStr][timeStr]; exists {
						takeId := new(uint) // 새로운 메모리 공간을 할당
						*takeId = val       // 값을 복사
						timeTaken[timeStr] = takeId
					} else {
						timeTaken[timeStr] = nil
					}
				}

				response := ExpectExerciseResponse{
					Id:        exercise.ID,
					Name:      exercise.Name,
					TimeTaken: timeTaken,
				}

				if _, exists := responseMap[dateStr]; !exists {
					responseMap[dateStr] = ExerciseTakeResponse{
						DateTaken:     dateStr,
						ExerciseTaken: []ExpectExerciseResponse{},
					}
				}

				// 맵에서 값을 가져와 변수에 저장
				tempResponse := responseMap[dateStr]
				// 변수의 필드를 수정
				tempResponse.ExerciseTaken = append(tempResponse.ExerciseTaken, response)
				// 수정된 변수를 다시 맵에 저장
				responseMap[dateStr] = tempResponse
			}
		}
	}

	// 맵을 리스트로 변환
	for _, response := range responseMap {
		responses = append(responses, response)
	}

	for _, take := range exerciseTakes {
		dateStr := take.DateTaken.Format("2006-01-02")
		timeStr := take.TimeTaken

		var found bool

		// 8.1. 해당 날짜의 기록이 있는지 확인합니다.
		for i, res := range responses {
			if res.DateTaken == dateStr {
				found = true
				var exerciseFound bool
				for j, medRes := range res.ExerciseTaken {
					if medRes.Id == take.ExerciseID {
						takeId := new(uint)
						*takeId = take.ID
						responses[i].ExerciseTaken[j].TimeTaken[timeStr] = takeId
						exerciseFound = true
						break
					}
				}
				if !exerciseFound {
					takeId := new(uint)
					*takeId = take.ID
					newExercise := ExpectExerciseResponse{
						Id:        take.ExerciseID,
						Name:      take.Exercise.Name,
						TimeTaken: map[string]*uint{timeStr: takeId},
					}
					responses[i].ExerciseTaken = append(responses[i].ExerciseTaken, newExercise)
				}
				break
			}
		}

		// 8.2. 해당 날짜에 대한 기록이 없으면 새로 추가합니다.
		if !found {
			takeId := new(uint)
			*takeId = take.ID
			timeTaken := map[string]*uint{
				timeStr: takeId,
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

func (service *exerciseService) getExercises(id uint) ([]ExerciseResponse, error) {
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

func (service *exerciseService) doExercise(request TakeExercise) (string, error) {

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

func (service *exerciseService) undoExercise(id, uid uint) (string, error) {

	if err := service.db.Where("id = ? AND uid = ?", id, uid).Unscoped().Delete(&model.ExerciseTake{}).Error; err != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}

// func (service *exerciseService) getProjects() ([]ProjectResponse, error) {

// 	var projects []ProjectResponse
// 	err := service.db.Model(&model.Video{}).
// 		Select("project_id, project_name as name, count(*) as count").
// 		Group("project_id").Scan(&projects).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return projects, nil
// }

// // face-service 의 face-exercise 쪽에는 한번에 다가져옴

// func (service *exerciseService) getVideos(projectId string, page uint) ([]VideoResponse, error) {
// 	pageSize := 20
// 	var videos []model.Video
// 	offset := page * uint(pageSize)
// 	if err := service.db.Where("project_id = ?", projectId).Offset(int(offset)).Limit(pageSize).Order("id DESC").Find(&videos).Error; err != nil {
// 		return nil, errors.New("db error")
// 	}

// 	var videoResponses []VideoResponse
// 	for _, v := range videos {
// 		videoResponses = append(videoResponses, VideoResponse{Name: v.Name, VideoId: v.VideoId, ThumbnailUrl: v.ThumbnailUrl, Duration: v.Duration})
// 	}

// 	return videoResponses, nil
// }
