package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func GetSampleVideosEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		videos, err := s.getSampleVideos()
		if err != nil {
			return nil, err
		}
		return videos, nil
	}
}

func GetFaceInfoEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetFaceInfoParams)
		scores, err := s.getFaceInfos(id, queryParams)
		if err != nil {
			return nil, err
		}
		return scores, nil
	}
}

func GetTapBlinkScoreEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetTapBlinkScoreParams)
		scores, err := s.getTapBlinkScores(id, queryParams)
		if err != nil {
			return nil, err
		}
		return scores, nil
	}
}

func SaveTapBlinkScoreEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		score := request.(TapBlinkRequest)
		code, err := s.saveTapBlinkScore(score)
		if err != nil {
			return nil, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SaveFaceInfoEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		score := reqMap["queryParams"].(FaceInfoRequest)
		code, err := s.saveFaceInfos(id, score)
		if err != nil {
			return nil, err
		}
		return BasicResponse{Code: code}, nil
	}
}
