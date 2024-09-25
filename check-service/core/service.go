package core

import (
	"check-service/model"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CheckService interface {
	getSampleVideos() ([]SampleVideoResponse, error)
	getFaceScores(id uint, param GetFaceScoreParams) ([]FaceScoreResponse, error)
	getTapBlinkScores(id uint, param GetTapBlinkScoreParams) ([]TapBlinkResponse, error)
	saveFaceScores(request FaceScoreRequest) (string, error)
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

func (service *checkService) getFaceScores(id uint, param GetFaceScoreParams) ([]FaceScoreResponse, error) {
	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}

	var faceScores []model.FaceScore
	var faceScoreResponses []FaceScoreResponse

	query := service.db.Where("uid = ?", id)
	if param.StartDate != "" {
		query = query.Where("created >= ?", param.StartDate)
	}
	if param.EndDate != "" {
		query = query.Where("created <= ?", param.EndDate+" 23:59:59")
	}
	query = query.Order("id DESC")

	if err := query.Find(&faceScores).Error; err != nil {
		return nil, errors.New("db error")
	}

	for _, v := range faceScores {
		faceScoreResponses = append(faceScoreResponses, FaceScoreResponse{
			TargetDate: v.CreatedAt.Format("2006-01-02"),
			FaceType:   v.FaceType,
			FaceLine:   v.FaceLine,
			Sd:         v.Sd,
		})
	}
	return faceScoreResponses, nil
}

func (service *checkService) getTapBlinkScores(id uint, param GetTapBlinkScoreParams) ([]TapBlinkResponse, error) {

	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}

	var tapBlinkScores []model.TapBlinkScore
	var tapBlinkResponses []TapBlinkResponse

	query := service.db.Where("uid = ?", id)
	if param.StartDate != "" {
		query = query.Where("created >= ?", param.StartDate)
	}
	if param.EndDate != "" {
		query = query.Where("created <= ?", param.EndDate+" 23:59:59")
	}
	query = query.Order("id DESC")

	if err := query.Find(&tapBlinkScores).Error; err != nil {
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

func (service *checkService) saveFaceScores(request FaceScoreRequest) (string, error) {
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

	if err := tx.Where("created_at::date = ? AND uid = ? ", targetDate, request.Uid).Delete(&model.FaceScore{}).Error; err != nil {
		tx.Rollback()
		return "", errors.New("db error")
	}

	var faceScore model.FaceScore
	faceScore.Uid = request.Uid
	faceScore.FaceLine = request.FaceLine
	faceScore.FaceType = request.FaceType
	faceScore.Sd = request.Sd

	if err := tx.Create(&faceScore).Error; err != nil {
		tx.Rollback()
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

	if err := tx.Where("created_at::date = ? AND uid = ? ", targetDate, request.Uid).Delete(&model.TapBlinkScore{}).Error; err != nil {
		tx.Rollback()
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
		return "", errors.New("db error2")
	}
	tx.Commit()
	return "200", nil
}
