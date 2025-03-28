package core

import (
	"emotion-service/model"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type EmotionService interface {
	saveEmotion(request EmotionRequest) (string, error)
	getEmotions(id uint, param GetEmotionsParams) ([]EmotionResponse, error)
}

type emotionService struct {
	db *gorm.DB
}

func NewEmotionService(db *gorm.DB) EmotionService {
	return &emotionService{db: db}
}

func (service *emotionService) saveEmotion(request EmotionRequest) (string, error) {
	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	if request.Symptoms != nil && len(request.Symptoms) != 0 {
		if len(request.Symptoms) > 5 {
			log.Println("invalid1")
			return "", errors.New("invalid Symptoms")
		}
		for _, v := range request.Symptoms {
			if v > 5 || v == 0 {
				log.Println("invalid2")
				return "", errors.New("invalid Symptoms")
			}
		}
	}

	targetDate, err := time.Parse("2006-01-02", request.TargetDate)
	if err != nil {
		log.Println(err)
		return "", errors.New("invalid TargetDate")
	}

	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	var emotion model.Emotion
	if err := tx.Where("target_date = ? AND uid = ? ", targetDate, request.Uid).Delete(&model.Emotion{}).Error; err != nil {
		tx.Rollback()
		log.Println("db error")
		return "", errors.New("db error")
	}

	emotion.Uid = request.Uid
	emotion.Emotion = request.Emotion
	emotion.Symptoms = uintSliceToInt64Array(request.Symptoms)
	emotion.Memo = request.Memo
	emotion.TargetDate = targetDate

	if err := tx.Create(&emotion).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}
	tx.Commit()
	return "200", nil
}

func (service *emotionService) getEmotions(id uint, param GetEmotionsParams) ([]EmotionResponse, error) {

	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}
	var emotions []model.Emotion

	if err := service.db.Where("uid = ? AND created_at >= ? AND created_at <= ? ", id, param.StartDate, param.EndDate+" 23:59:59").Find(&emotions).Order("id DESC").Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	var emotionResponses []EmotionResponse
	for _, v := range emotions {
		emotionResponses = append(emotionResponses, EmotionResponse{Emotion: v.Emotion, Symptoms: int64ArrayToUintSlice(v.Symptoms), Memo: v.Memo, TargetDate: v.TargetDate.Format("2006-01-02")})
	}

	return emotionResponses, nil
}
