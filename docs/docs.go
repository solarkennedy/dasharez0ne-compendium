// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Shows all macros by keyword. Searches captions, original_text, and tags",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Search macros",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Keyword",
                        "name": "keyword",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Macro": {
            "type": "object",
            "properties": {
                "caption": {
                    "type": "string"
                },
                "edit_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "original_text": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
