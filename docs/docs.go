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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "logs in as an exiting users using login and passowrd",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "log in",
                "parameters": [
                    {
                        "description": "login details",
                        "name": "userDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.logInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.authResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "403": {
                        "description": "forbidden (invalid password)",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not exists",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "creates a new user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "sign up as a new user",
                "parameters": [
                    {
                        "description": "sign up details",
                        "name": "userDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.signUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.authResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "409": {
                        "description": "conflict (user with such a name exists)",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    }
                }
            }
        },
        "/location/fetch": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "fetches users location",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "location"
                ],
                "summary": "returns location of provided users",
                "parameters": [
                    {
                        "description": "arrays of ids or nicks",
                        "name": "fetchLocation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/location.fetchRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "list of user",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/location.userLocation"
                            }
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    }
                }
            }
        },
        "/location/update": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "updates user's location",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "location"
                ],
                "summary": "update user's location",
                "parameters": [
                    {
                        "description": "location",
                        "name": "locationUpdate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/location.updateLocationRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "invalid request",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/jsonerr.JSONError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.authResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID is user id (uuid)",
                    "type": "string"
                },
                "refresh_token": {
                    "description": "RefreshToken user refresh token",
                    "type": "string"
                },
                "token": {
                    "description": "Token user auth token (Bearer)",
                    "type": "string"
                }
            }
        },
        "auth.logInRequest": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "description": "Name username",
                    "type": "string"
                },
                "password": {
                    "description": "Password user password",
                    "type": "string"
                }
            }
        },
        "auth.signUpRequest": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "description": "Name username, must be unique",
                    "type": "string"
                },
                "password": {
                    "description": "Password user password, min 8 characters",
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "jsonerr.JSONError": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Code is desired http code for this error",
                    "type": "integer"
                },
                "error": {
                    "description": "Err is a golang error returned by the app\nIt is removed in production application (TBD)",
                    "type": "string"
                },
                "message": {
                    "description": "Message is human friendly error message",
                    "type": "string"
                }
            }
        },
        "location.fetchRequest": {
            "type": "object",
            "properties": {
                "nicks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "uuids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "location.updateLocationRequest": {
            "type": "object",
            "required": [
                "latitude",
                "longitude"
            ],
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "location.userLocation": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "last_update": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "nick": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "WhereIsEveryone",
	Description:      "This is a sample server for WhereIsEveryone",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
