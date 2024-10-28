package core

import (
	"check-service/model"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CheckService interface {
	getSampleVideos() ([]SampleVideoResponse, error)
	getFaceInfos(id uint, param GetFaceInfoParams) ([]FaceInfoResponse, error)
	getTapBlinkScores(id uint, param GetTapBlinkScoreParams) ([]TapBlinkResponse, error)
	saveFaceInfos(uid uint, request FaceInfoRequest) (string, error)
	saveTapBlinkScore(request TapBlinkRequest) (string, error)
}

type checkService struct {
	db *gorm.DB
}

func NewCheckService(db *gorm.DB) CheckService {
	return &checkService{db: db}
}

func (service *checkService) getSampleVideos() ([]SampleVideoResponse, error) {
	var sampleVideos []model.SampleVideo
	var sampleVideoResponses []SampleVideoResponse

	if err := service.db.Find(&sampleVideos).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	for _, v := range sampleVideos {
		sampleVideoResponses = append(sampleVideoResponses, SampleVideoResponse{
			Category:  v.Category,
			VideoType: v.VideoType,
			Title:     v.Titile,
			VideoId:   v.VideoId,
		})
	}

	return sampleVideoResponses, nil
}

func (service *checkService) getFaceInfos(id uint, param GetFaceInfoParams) ([]FaceInfoResponse, error) {
	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}

	var faceInfos []model.FaceInfo

	query := service.db.Where("uid = ?", id)
	if param.StartDate != "" {
		query = query.Where("created_at >= ?", param.StartDate)
	}
	if param.EndDate != "" {
		query = query.Where("created_at <= ?", param.EndDate+" 23:59:59")
	}

	if err := query.Find(&faceInfos).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	// 2. 날짜별로 그룹화된 응답을 저장할 맵
	responseMap := make(map[string]FaceInfoResponse)

	// 3. 데이터를 날짜 및 FaceType별로 그룹화하여 맵에 저장
	for _, v := range faceInfos {
		dateStr := v.CreatedAt.Format("2006-01-02")

		// 3.1 해당 날짜에 대한 응답이 없으면 새로 생성
		if _, exists := responseMap[dateStr]; !exists {
			responseMap[dateStr] = FaceInfoResponse{
				TargetDate: dateStr,
				FaceInfos:  make(map[uint][]FaceInfo), // FaceType을 키로 설정
			}
		}

		// 3.2 FaceInfo 생성
		faceInfo := FaceInfo{
			FaceLine: v.FaceLine,
			Sd:       v.Sd,
		}

		// 3.3 해당 날짜의 응답을 가져와서 FaceType별로 추가
		response := responseMap[dateStr]
		response.FaceInfos[v.FaceType] = append(response.FaceInfos[v.FaceType], faceInfo)
		responseMap[dateStr] = response
	}

	// 4. 맵을 리스트로 변환
	var responses []FaceInfoResponse
	for _, response := range responseMap {
		responses = append(responses, response)
	}

	// 5. responses를 날짜 순으로 정렬
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].TargetDate < responses[j].TargetDate
	})

	return responses, nil
}

func (service *checkService) getTapBlinkScores(id uint, param GetTapBlinkScoreParams) ([]TapBlinkResponse, error) {

	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}

	var tapBlinkScores []model.TapBlinkScore
	var tapBlinkResponses []TapBlinkResponse

	query := service.db.Where("uid = ? AND score_type= ? ", id, param.ScoreType)
	if param.StartDate != "" {
		query = query.Where("created_at >= ?", param.StartDate)
	}
	if param.EndDate != "" {
		query = query.Where("created_at <= ?", param.EndDate+" 23:59:59")
	}
	query = query.Order("id DESC")

	if err := query.Find(&tapBlinkScores).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	for _, v := range tapBlinkScores {
		tapBlinkResponses = append(tapBlinkResponses, TapBlinkResponse{
			TargetDate:    v.CreatedAt.Format("2006-01-02"),
			SuccessCount:  v.SuccessCount,
			ErrorCount:    v.ErrorCount,
			ScoreType:     v.ScoreType,
			ReactionSpeed: v.ReactionSpeed,
		})
	}
	return tapBlinkResponses, nil
}

func (service *checkService) saveFaceInfos(uid uint, request FaceInfoRequest) (string, error) {
	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	targetDate := time.Now().Format("2006-01-02")

	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// 1. 삭제할 faceType 목록을 추출
	var faceTypes []uint
	for faceType := range request.FaceInfos {
		if faceType > 5 || faceType == 0 {
			tx.Rollback()
			log.Println("check")
			return "", errors.New("check face_type")
		}
		faceTypes = append(faceTypes, faceType)
	}

	// 2. 동일한 uid, targetDate, face_type에 해당하는 기존 데이터를 삭제 (IN 절 사용)
	if len(faceTypes) > 0 {
		if err := tx.Where("created_at::date = ? AND uid = ? AND face_type IN (?)", targetDate, uid, faceTypes).Unscoped().Delete(&model.FaceInfo{}).Error; err != nil {
			tx.Rollback()
			log.Println("db error")
			return "", errors.New("db error")
		}
	}

	// 3. 데이터를 저장할 슬라이스 생성
	var models []model.FaceInfo

	// 4. request.FaceInfos에서 데이터를 추출하여 저장할 준비
	for faceType, faceInfos := range request.FaceInfos {
		for _, faceInfo := range faceInfos {
			models = append(models, model.FaceInfo{
				Uid:      uid,
				FaceType: faceType,
				FaceLine: faceInfo.FaceLine,
				Sd:       faceInfo.Sd,
			})
		}
	}

	if err := tx.Create(&models).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}
	tx.Commit()
	return "200", nil
}

func (service *checkService) saveTapBlinkScore(request TapBlinkRequest) (string, error) {
	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		log.Println(err)
		return "", err
	}

	targetDate := time.Now().Format("2006-01-02")

	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if err := tx.Where("created_at::date = ? AND uid = ? AND score_type=? ", targetDate, request.Uid, request.ScoreType).Unscoped().Delete(&model.TapBlinkScore{}).Error; err != nil {
		tx.Rollback()
		log.Println("db error")
		return "", errors.New("db error")
	}

	var tapBlink model.TapBlinkScore
	tapBlink.Uid = request.Uid
	tapBlink.ScoreType = request.ScoreType
	tapBlink.ErrorCount = request.ErrorCount
	tapBlink.ReactionSpeed = request.ReactionSpeed
	tapBlink.SuccessCount = request.SuccessCount
	tapBlink.ReactionSpeed = request.ReactionSpeed

	if err := tx.Create(&tapBlink).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}
	tx.Commit()
	return "200", nil
}
