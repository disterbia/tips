package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func SaveExerciseEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(ExerciseRequest)
		code, err := s.saveExercise(exercise)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetExpectsEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		inquires, err := s.getExpects(uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return inquires, nil
	}
}

func GetExercisesEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		medicines, err := s.getExercises(id)
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
		code, err := s.removeExercise(id, uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func DoExerciseEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		exercise := request.(TakeExercise)
		code, err := s.doExercise(exercise)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetProjectsEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		projects, err := s.getProjects()
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return projects, nil
	}
}

func GetVideosEndpoint(s ExerciseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(GetVideoParams)
		videos, err := s.getVideos(reqMap.ProjectId, reqMap.Page)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return videos, nil
	}
}
