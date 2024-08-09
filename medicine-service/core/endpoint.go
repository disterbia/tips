package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		medicine := request.(MedicineRequest)
		code, err := s.SaveMedicine(medicine)
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
		code, err := s.RemoveMedicines(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetExpectsEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		inquires, err := s.GetExpects(uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetMedicinesEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		medicines, err := s.GetMedicines(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}

func TakeEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		takeMedicine := request.(TakeMedicine)
		code, err := s.TakeMedicine(takeMedicine)
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
		code, err := s.UnTakeMedicine(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SearchsEndpoint(s MedicineService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		keyword := request.(string)
		medicines, err := s.SearchMedicines(keyword)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}
