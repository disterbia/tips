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
        "/do-exercises": {
            "post": {
                "description": "운동 완료 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 기록",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "운동 완료 데이터",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.TakeExercise"
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
        "/get-exercises": {
            "get": {
                "description": "등록 운동 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "등록 운동 조회",
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
                        "description": "등록 운동 정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.ExerciseResponse"
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
        "/get-projects": {
            "get": {
                "description": "운동 동영상 카테고리 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 동영상 카테고리 조회",
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
                        "description": "카테고리 정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.ProjectResponse"
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
                "description": "운동 내역 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 내역 조회",
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
                        "description": "운동 내역정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.ExerciseTakeResponse"
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
        "/get-videos": {
            "get": {
                "description": "카테고리별 운동 동영상 조회시 호출",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "카테고리별 운동 동영상 조회 (20개씩)",
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
                        "description": "project_id",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "페이지 default 0",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "동영상 정보",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.VideoResponse"
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
        "/remove-exercise/{id}": {
            "post": {
                "description": "운동 삭제시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 삭제",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "삭제할 id 배열",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
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
        "/save-exercise": {
            "post": {
                "description": "운동 등록 및 수정시 호출 - 생성시 id생략",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 저장",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer {jwt_token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "요청 DTO - 운동 데이터",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/core.ExerciseRequest"
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
        "/undo-exercise/{id}": {
            "post": {
                "description": "운동 취소시 호출",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "운동 /exercise"
                ],
                "summary": "운동 기록",
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
                        "description": "취소 ID",
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
        "core.ExerciseRequest": {
            "type": "object",
            "properties": {
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
                "name": {
                    "type": "string"
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
                "weekdays": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.ExerciseResponse": {
            "type": "object",
            "properties": {
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
                "name": {
                    "type": "string"
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
                "weekdays": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.ExerciseTakeResponse": {
            "type": "object",
            "properties": {
                "Exercise_taken": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.ExpectExerciseResponse"
                    }
                },
                "date_taken": {
                    "type": "string",
                    "example": "YYYY-MM-dd"
                }
            }
        },
        "core.ExpectExerciseResponse": {
            "type": "object",
            "properties": {
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
        "core.ProjectResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "project_id": {
                    "type": "string"
                }
            }
        },
        "core.TakeExercise": {
            "type": "object",
            "properties": {
                "date_taken": {
                    "type": "string",
                    "example": "YYYY-MM-DD"
                },
                "exercise_id": {
                    "type": "integer"
                },
                "time_taken": {
                    "type": "string",
                    "example": "HH:mm"
                }
            }
        },
        "core.VideoResponse": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "thumbnail_url": {
                    "type": "string"
                },
                "video_id": {
                    "type": "string"
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
