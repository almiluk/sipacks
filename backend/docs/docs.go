// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "almiluk",
            "email": "almiluk@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/packs": {
            "get": {
                "description": "List packs with filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "packs"
                ],
                "summary": "List packs",
                "parameters": [
                    {
                        "description": "Filter",
                        "name": "filter",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PackListRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PackListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add new questions pack",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "packs"
                ],
                "summary": "Add pack",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Pack data",
                        "name": "pack",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AddPackResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/packs/{id}": {
            "get": {
                "description": "Download questions pack",
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "packs"
                ],
                "summary": "Download pack",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Pack ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AddPackResponse": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string",
                    "example": "author"
                },
                "creation_date": {
                    "type": "string",
                    "example": "creation_date"
                },
                "name": {
                    "type": "string",
                    "example": "name"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "error"
                },
                "message": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "models.PackListRequest": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string",
                    "example": "author"
                },
                "max_creation_date": {
                    "type": "string",
                    "example": "01.01.1970"
                },
                "min_creation_date": {
                    "type": "string",
                    "example": "01.01.1970"
                },
                "name": {
                    "type": "string",
                    "example": "name"
                },
                "sort_by": {
                    "type": "string",
                    "enum": [
                        "creation_date",
                        "downloads_num"
                    ],
                    "example": "creation_date"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "tags"
                    ]
                }
            }
        },
        "models.PackListResponse": {
            "type": "object",
            "properties": {
                "packs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PackResponse"
                    }
                },
                "packs_num": {
                    "type": "integer",
                    "example": 0
                }
            }
        },
        "models.PackResponse": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string",
                    "example": "author"
                },
                "creation_date": {
                    "type": "string",
                    "example": "creation_date"
                },
                "downloads_num": {
                    "type": "integer",
                    "example": 0
                },
                "file_size": {
                    "type": "integer",
                    "example": 0
                },
                "name": {
                    "type": "string",
                    "example": "name"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API for managing question packs for the 'SIGame' game",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
