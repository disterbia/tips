package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func AdminLoginEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		email := reqMap["email"].(string)
		password := reqMap["password"].(string)

		token, err := s.AdminLogin(email, password)

		if err != nil {
			return LoginResponse{Err: err.Error()}, err
		}
		return LoginResponse{Jwt: token}, nil
	}
}

func AnswerEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(InquireReplyRequest)
		code, err := s.AnswerInquire(answer)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(InquireRequest)
		code, err := s.SendInquire(inquire)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func RemoveInquireEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveInquire(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}
func RemoveReplyEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveReply(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetInquireParams)
		inquires, err := s.GetMyInquires(id, queryParams.Page, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetAllEndpoint(s InquireService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		queryParams := reqMap["queryParams"].(GetInquireParams)
		inquires, err := s.GetAllInquires(id, queryParams.Page, queryParams.StartDate, queryParams.EndDate)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}
