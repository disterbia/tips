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

func GetFaceScoreEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetFaceScoreParams)
		scores, err := s.getFaceScores(id, queryParams)
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

func SaveFaceScoreEndpoint(s CheckService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		score := reqMap["queryParams"].([]FaceScoreRequest)
		code, err := s.saveFaceScores(id, score)
		if err != nil {
			return nil, err
		}
		return BasicResponse{Code: code}, nil
	}
}
