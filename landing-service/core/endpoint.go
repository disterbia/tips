package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func KldgaInquireEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(KldgaInquireRequest)
		code, err := s.kldgaInquire(inquire)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func KldgaCompetitionEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		competition := request.(KldgaCompetitionRequest)
		code, err := s.kldgaCompetition(competition)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func AdapfitInqruieEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(AdapfitInquireReqeust)
		code, err := s.adapfitInquire(inquire)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

// 인증번호 발송 엔드포인트
func SendAuthCodeEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthCodeRequest)
		code, err := s.sendAuthCode(req.Phone)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

// 인증번호 검증 엔드포인트
func VerifyAuthCodeEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyAuthRequest)
		code, err := s.verifyAuthCode(req.Phone, req.Code)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}
