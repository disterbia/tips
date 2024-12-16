package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func KldgaInquireEndpoint(s LandingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		inquire := request.(KldgaRequest)
		code, err := s.kldgaInquire(inquire)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}
