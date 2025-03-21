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
        "/get-medicines": {
            "get": {
                "description": "등록 약물 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "등록 약물 조회",
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
                        "description": "등록 약물 정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.MedicineResponse"
                            }
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
        "/get-takens": {
            "get": {
                "description": "약물 복용내역 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 복용내역 조회",
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
                        "description": "운동정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.MedicineTakeResponse"
                            }
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
        "/remove-medicine/{id}": {
            "post": {
                "description": "약물 삭제시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 삭제",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "약물 ID",
                        "name": "id",
                        "in": "path"
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
        "/save-medicine": {
            "post": {
                "description": "약물등록 및 수정시 호출 - 생성시 id생략",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 저장",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "요청 DTO - 약물데이터",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.MedicineRequest"
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
        "/search-medicines": {
            "get": {
                "description": "약물 검색 키워드 입력시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 찾기",
                "parameters": [
                    {
                        "type": "string",
                        "description": "키워드",
                        "name": "keyword",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "약물명",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
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
        "/take-medicine": {
            "post": {
                "description": "약물 복용시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 복용",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "약물 복용 데이터",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.TakeMedicine"
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
        "/untake-medicine/{id}": {
            "post": {
                "description": "약물 복용취소시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "약물 /medicine"
                ],
                "summary": "약물 복용취소",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "복용 ID",
                        "name": "id",
                        "in": "path"
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
        }
    },
    "definitions": {
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
        "core.ExpectMedicineResponse": {
            "type": "object",
            "properties": {
                "dose": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "time_taken": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.MedicineRequest": {
            "type": "object",
            "properties": {
                "dose": {
                    "type": "number"
                },
                "end_at": {
                    "type": "string",
                    "example": "YYYY-MM:dd"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "medicine_type": {
                    "type": "string"
                },
                "min_reserves": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "remaining": {
                    "type": "number"
                },
                "start_at": {
                    "type": "string",
                    "example": "YYYY-MM-dd"
                },
                "times": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "HH:mm",
                        "HH:mm"
                    ]
                },
                "use_privacy": {
                    "type": "boolean"
                },
                "weekdays": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.MedicineResponse": {
            "type": "object",
            "properties": {
                "dose": {
                    "type": "number"
                },
                "end_at": {
                    "type": "string",
                    "example": "YYYY-MM:dd"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "medicine_type": {
                    "type": "string"
                },
                "min_reserves": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "remaining": {
                    "type": "number"
                },
                "start_at": {
                    "type": "string",
                    "example": "YYYY-MM-dd"
                },
                "times": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "use_privacy": {
                    "type": "boolean"
                },
                "weekdays": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.MedicineTakeResponse": {
            "type": "object",
            "properties": {
                "date_taken": {
                    "type": "string",
                    "example": "YYYY-MM-dd"
                },
                "medicine_taken": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.ExpectMedicineResponse"
                    }
                }
            }
        },
        "core.TakeMedicine": {
            "type": "object",
            "properties": {
                "date_taken": {
                    "type": "string",
                    "example": "YYYY-MM-DD"
                },
                "dose": {
                    "type": "number"
                },
                "medicine_id": {
                    "type": "integer"
                },
                "time_taken": {
                    "type": "string",
                    "example": "HH:mm"
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
