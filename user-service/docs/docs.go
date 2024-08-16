// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auto-login": {
            "post": {
                "security": [
                    {
                        "jwt": []
                    }
                ],
                "description": "최초 로그인 이후 앱 실행시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "로그인 /user"
                ],
                "summary": "자동로그인",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "요청 DTO",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.AutoLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 JWT 토큰 반환",
                        "schema": {
                            "$ref": "#/definitions/core.SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/get-user": {
            "get": {
                "description": "내 정보 조회시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "회원정보 조회 /user"
                ],
                "summary": "회원정보 조회",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 유저 객체 반환/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자 sns_type- 0:휴대폰,1:카카오 2:구글 3:애플",
                        "schema": {
                            "$ref": "#/definitions/core.UserResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/get-version": {
            "get": {
                "description": "최신버전 조회시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "공통 /user"
                ],
                "summary": "최신버전 조회",
                "responses": {
                    "200": {
                        "description": "최신 버전 정보",
                        "schema": {
                            "$ref": "#/definitions/core.AppVersionResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/link-email": {
            "post": {
                "description": "계정 연동시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "계정 연동 /user"
                ],
                "summary": "계정 연동",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "요청 DTO",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.LinkRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/phone-login": {
            "post": {
                "description": "휴대번호 로그인시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "로그인 /user"
                ],
                "summary": "휴대번호 로그인",
                "parameters": [
                    {
                        "description": "요청 DTO user- user_type: 0:해당없음, 1~6:파킨슨 환자, 10:보호자 / 최초 로그인 이후 로그인시 phone,fcm_token,device_id 만 필요함",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.PhoneLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 JWT 토큰 반환",
                        "schema": {
                            "$ref": "#/definitions/core.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환: 오류메시지 KAKAO=1, GOOGLE=2, APPLE=3 / '-1' = 인증필요 , '-2' = 추가정보 입력 필요 ",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/remove-user": {
            "post": {
                "description": "회원탈퇴시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "회원탈퇴 /user"
                ],
                "summary": "회원탈퇴",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/send-code-join/{number}": {
            "post": {
                "description": "회원가입 인증번호 발송시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "인증번호 /user"
                ],
                "summary": "인증번호 발송",
                "parameters": [
                    {
                        "type": "string",
                        "description": "휴대번호",
                        "name": "number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환: 오류메시지 \"-1\" = 이미 가입한번호",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/send-code-login/{number}": {
            "post": {
                "description": "휴대번호 로그인 인증번호 발송시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "인증번호 /user"
                ],
                "summary": "인증번호 발송",
                "parameters": [
                    {
                        "type": "string",
                        "description": "휴대번호",
                        "name": "number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sns-login": {
            "post": {
                "responses": {}
            }
        },
        "/update-user": {
            "post": {
                "description": "내정보 변경시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "마이페이지 /user"
                ],
                "summary": "내정보 변경",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "요청 DTO - 업데이트 할 데이터/ ture:남성 user_Type- 0:해당없음 1:파킨슨 환자 2:보호자",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환: 오류메시지 \"-1\" = 번호인증 필요",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/verify-code": {
            "post": {
                "description": "인증번호 입력 후 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "인증번호 /user"
                ],
                "summary": "번호 인증",
                "parameters": [
                    {
                        "description": "요청 DTO",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.VerifyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공시 200 반환",
                        "schema": {
                            "$ref": "#/definitions/core.BasicResponse"
                        }
                    },
                    "400": {
                        "description": "요청 처리 실패시 오류 메시지 반환",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "요청 처리 실패시 오류 메시지 반환: 오류메시지 \"-1\" = 코드불일치",
                        "schema": {
                            "$ref": "#/definitions/core.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "core.AppVersionResponse": {
            "type": "object",
            "properties": {
                "android_link": {
                    "type": "string"
                },
                "ios_link": {
                    "type": "string"
                },
                "latest_version": {
                    "type": "string"
                }
            }
        },
        "core.AutoLoginRequest": {
            "type": "object",
            "properties": {
                "device_id": {
                    "type": "string"
                },
                "fcm_token": {
                    "type": "string"
                }
            }
        },
        "core.BasicResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "core.ErrorResponse": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "core.ImageResponse": {
            "type": "object",
            "properties": {
                "thumbnail_url": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "core.LinkRequest": {
            "type": "object",
            "properties": {
                "id_token": {
                    "type": "string"
                }
            }
        },
        "core.LinkedResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "sns_type": {
                    "type": "integer"
                }
            }
        },
        "core.PhoneLoginRequest": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "string",
                    "example": "yyyy-mm-dd"
                },
                "device_id": {
                    "type": "string"
                },
                "fcm_token": {
                    "type": "string"
                },
                "gender": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "user_type": {
                    "type": "integer"
                }
            }
        },
        "core.SuccessResponse": {
            "type": "object",
            "properties": {
                "jwt": {
                    "type": "string"
                }
            }
        },
        "core.UserRequest": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "string",
                    "example": "yyyy-mm-dd"
                },
                "gender": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "profile_image": {
                    "type": "string",
                    "example": "base64string"
                },
                "user_type": {
                    "type": "integer"
                }
            }
        },
        "core.UserResponse": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "string",
                    "example": "yyyy-mm-dd"
                },
                "created_at": {
                    "type": "string"
                },
                "gender": {
                    "description": "true:남 false: 여",
                    "type": "boolean"
                },
                "linked_emails": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.LinkedResponse"
                    }
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "profile_image": {
                    "$ref": "#/definitions/core.ImageResponse"
                },
                "sns_type": {
                    "type": "integer"
                }
            }
        },
        "core.VerifyRequest": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "인증번호 6자리"
                },
                "phone_number": {
                    "type": "string",
                    "example": "01000000000"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
