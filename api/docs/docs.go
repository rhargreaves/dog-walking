// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/dogs": {
            "get": {
                "description": "Get a list of all registered dogs",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dogs"
                ],
                "summary": "List all dogs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Dog"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new dog with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dogs"
                ],
                "summary": "Create a new dog",
                "parameters": [
                    {
                        "description": "Dog information",
                        "name": "dog",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Dog"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Dog"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            }
        },
        "/dogs/{id}": {
            "get": {
                "description": "Get details of a specific dog by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dogs"
                ],
                "summary": "Get a dog by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Dog"
                        }
                    },
                    "404": {
                        "description": "Dog not found",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a dog's information by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dogs"
                ],
                "summary": "Update a dog",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated dog information",
                        "name": "dog",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Dog"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Dog"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "404": {
                        "description": "Dog not found",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a dog by its ID",
                "tags": [
                    "dogs"
                ],
                "summary": "Delete a dog",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Dog not found",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            }
        },
        "/dogs/{id}/photo": {
            "put": {
                "description": "Upload a JPEG photo for a specific dog",
                "consumes": [
                    "image/jpeg"
                ],
                "tags": [
                    "dogs",
                    "photos"
                ],
                "summary": "Upload a dog's photo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Invalid content type or request",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "404": {
                        "description": "Dog not found",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            }
        },
        "/dogs/{id}/photo/detect-breed": {
            "post": {
                "description": "Analyzes a previously uploaded photo to detect the dog's breed",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dogs",
                    "photos"
                ],
                "summary": "Detect a dog's breed from its photo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Dog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns id, breed, and confidence",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "No dog detected or no specific breed detected",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "404": {
                        "description": "Dog not found",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/common.APIError"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Returns OK if the API is running",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Dog": {
            "type": "object",
            "properties": {
                "breed": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "photoHash": {
                    "type": "string"
                },
                "photoUrl": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api.dog-walking.com",
	BasePath:         "/",
	Schemes:          []string{"https"},
	Title:            "Dog Walking API",
	Description:      "API for managing dogs in a dog walking service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
