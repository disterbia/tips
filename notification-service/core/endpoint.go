package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func GetMessagesEndpoint(s NotificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		notis, err := s.GetMessages(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return notis, nil
	}
}

func ReadAllEndpoint(s NotificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		code, err := s.ReadAll(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func RemoveMessagesEndpoint(s NotificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		ids := reqMap["ids"].([]uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveMessages(ids, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}
