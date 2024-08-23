package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func LoginEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		info := request.(LoginRequest)
		token, err := s.login(info)
		if err != nil {
			return nil, err
		}
		return LoginResponse{Jwt: token}, nil
	}
}

func SearchHospitalsEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		param := request.(SearchParam)
		result, err := s.searchHospitals(param)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func GetPoliciesEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := s.getPolicies()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
func VerifyEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		veri := request.(VerifyRequest)
		code, err := s.verifyCode(veri.PhoneNumber, veri.Code)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendCodeForSignInEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		number := request.(string)
		code, err := s.sendAuthCodeForSignin(number)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendCodeForIdEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindIdRequest)
		code, err := s.sendAuthCodeForId(req)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendCodeForPwEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindPwRequest)
		code, err := s.sendAuthCodeForPw(req)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func ChangePwEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindPasswordRequest)
		code, err := s.changePw(req)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func FindIdEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindIdRequest)
		email, err := s.findId(req)
		if err != nil {
			return "", err
		}
		return email, nil
	}
}

func SignInEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		param := request.(SignInRequest)
		result, err := s.signIn(param)
		if err != nil {
			return nil, err
		}
		return BasicResponse{Code: result}, nil
	}
}

func QuestionEndpoint(s AdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		param := request.(QuestionRequest)
		result, err := s.question(param)
		if err != nil {
			return nil, err
		}
		return BasicResponse{Code: result}, nil
	}
}
