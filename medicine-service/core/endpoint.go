package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		medicine := request.(MedicineRequest)
		code, err := s.saveMedicine(medicine)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func RemoveEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.removeMedicine(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetExpectsEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		inquires, err := s.getExpects(uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetMedicinesEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		medicines, err := s.getMedicines(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}

func TakeEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		takeMedicine := request.(TakeMedicine)
		code, err := s.takeMedicine(takeMedicine)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func UnTakeEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.unTakeMedicine(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SearchsEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		keyword := request.(string)
		medicines, err := s.searchMedicines(keyword)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}
