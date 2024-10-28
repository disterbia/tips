package core

import (
	"context"
	"errors"
	"log"

	"github.com/go-kit/kit/endpoint"
)

func SnsLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := s.snsLogin(req)
		if err != nil {
			return LoginResponse{Err: err.Error()}, err
		}
		return LoginResponse{Jwt: token}, nil
	}
}

func PhoneLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PhoneLoginRequest)
		token, err := s.phoneLogin(req)
		if err != nil {
			return LoginResponse{Err: err.Error()}, err
		}
		return LoginResponse{Jwt: token}, nil
	}
}

func AutoLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		autoRequest := request.(AutoLoginRequest)
		token, err := s.autoLogin(autoRequest)
		if err != nil {
			return LoginResponse{Err: err.Error()}, err
		}
		return LoginResponse{Jwt: token}, nil
	}
}

func VerifyEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		veri := request.(VerifyRequest)
		code, err := s.verifyAuthCode(veri.PhoneNumber, veri.Code)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendCodeForSignInEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		number := request.(string)
		code, err := s.sendAuthCodeForSingin(number)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func SendCodeForLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		number := request.(string)
		code, err := s.sendAuthCodeForLogin(number)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func UpdateUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := request.(UserRequest)
		code, err := s.updateUser(user)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}
func GetUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(uint)
		result, err := s.getUser(id)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return result, nil
	}
}

func LinkEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqMap := request.(LinkRequest)

		code, err := s.linkEmail(reqMap.Id, reqMap.IdToken)

		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func RemoveEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid := request.(uint)
		code, err := s.removeUser(uid)
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return BasicResponse{Code: code}, nil
	}
}

func GetVersionEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		version, err := s.getVersion()
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return version, nil
	}
}

func GetPolicesEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		version, err := s.getPolices()
		if err != nil {
			return BasicResponse{Code: err.Error()}, err
		}
		return version, nil
	}
}

func AppleCallbackEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AppleCallbackRequest)

		if req.Code == "" {
			log.Println("code")
			return nil, errors.New("authorization code is missing")
		}

		// Authorization Code로 애플과 통신해 토큰을 교환합니다.
		tokenResponse, err := s.exchangeCodeForToken(req.Code)
		if err != nil {
			return nil, err
		}

		return AppleCallbackResponse{
			AccessToken: tokenResponse.AccessToken,
			IDToken:     tokenResponse.IDToken,
		}, nil
	}
}
