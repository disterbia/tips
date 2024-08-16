package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveEmotionEndpoint(s EmotionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		emotion := request.(EmotionRequest)
		code, err := s.saveEmotion(emotion)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetEmotionsEndpoint(s EmotionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetEmotionsParams)
		inquires, err := s.getEmotions(id, queryParams)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}
