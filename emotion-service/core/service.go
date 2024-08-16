package core

import (
	"emotion-service/model"
	"errors"
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
			return "", errors.New("invalid Symptoms")
		}
		for _, v := range request.Symptoms {
			if v > 5 || v == 0 {
				return "", errors.New("invalid Symptoms")
			}
		}
	}

	targetDate, err := time.Parse("2006-01-02", request.TargetDate)
	if err != nil {
		return "", errors.New("invalid TargetDate")
	}

	var emotion model.Emotion
	if err := service.db.Where("id = ? AND uid = ? ", request.Id, request.Uid).First(&emotion).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("db error")
		}

	}

	emotion.Uid = request.Uid
	emotion.Emotion = request.Emotion
	emotion.Symptoms = uintSliceToInt64Array(request.Symptoms)
	emotion.Memo = request.Memo
	emotion.TargetDate = targetDate

	if err := service.db.Save(&emotion).Error; err != nil {
		return "", errors.New("db error2")
	}

	return "200", nil
}

func (service *emotionService) getEmotions(id uint, param GetEmotionsParams) ([]EmotionResponse, error) {

	validate := validator.New()

	if err := validate.Struct(param); err != nil {
		return nil, err
	}
	var emotions []model.Emotion

	if err := service.db.Where("uid = ? AND created_at >= ? AND created_at <= ? ", id, param.StartDate, param.EndDate+" 23:59:59").Find(&emotions).Order("id DESC").Error; err != nil {
		return nil, errors.New("db error")
	}

	var emotionResponses []EmotionResponse
	for _, v := range emotions {
		emotionResponses = append(emotionResponses, EmotionResponse{Id: v.ID, Emotion: v.Emotion, Symptoms: int64ArrayToUintSlice(v.Symptoms), Memo: v.Memo, TargetDate: v.TargetDate.Format("2006-01-02")})
	}

	return emotionResponses, nil
}
