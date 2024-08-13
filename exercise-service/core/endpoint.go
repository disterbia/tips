package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveExerciseEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(ExerciseRequest)
		code, err := s.SaveExercise(exercise)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetExpectsEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		inquires, err := s.GetExpects(uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetExercisesEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		medicines, err := s.GetExercises(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return medicines, nil
	}
}

func RemoveExerciseEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(map[string]interface{})
		id := reqMap["id"].(uint)
		uid := reqMap["uid"].(uint)
		code, err := s.RemoveExercise(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func DoExerciseEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(TakeExercise)
		code, err := s.DoExercise(exercise)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetProjectsEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		projects, err := s.GetProjects()
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return projects, nil
	}
}

func GetVideosEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(GetVideoParams)
		videos, err := s.GetVideos(reqMap.ProjectId, reqMap.Page)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return videos, nil
	}
}
