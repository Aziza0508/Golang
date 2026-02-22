package handlers

import (
	"net/http"
	"strings"
)

func SwaggerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "doc.json") {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(swaggerJSON))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerUIHTML))
	}
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Practice 3 API - Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger/doc.json",
            dom_id: '#swagger-ui',
        })
    </script>
</body>
</html>`

const swaggerJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Practice 3 - User Management API",
    "description": "A Go REST API with PostgreSQL, Redis caching, JWT auth, and more.",
    "version": "1.0.0"
  },
  "servers": [
    { "url": "http://localhost:8080", "description": "Local dev" }
  ],
  "security": [
    { "ApiKeyAuth": [] }
  ],
  "paths": {
    "/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Healthcheck",
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": { "application/json": { "schema": { "type": "object", "properties": { "status": { "type": "string", "example": "ok" } } } } }
          }
        }
      }
    },
    "/auth/register": {
      "post": {
        "tags": ["Auth"],
        "summary": "Register a new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateUserRequest" }
            }
          }
        },
        "responses": {
          "201": { "description": "User registered successfully" },
          "400": { "description": "Validation error" }
        }
      }
    },
    "/auth/login": {
      "post": {
        "tags": ["Auth"],
        "summary": "Login and get JWT token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password"],
                "properties": {
                  "email": { "type": "string", "example": "admin@example.com" },
                  "password": { "type": "string", "example": "secret123" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": { "application/json": { "schema": { "type": "object", "properties": { "token": { "type": "string" } } } } }
          },
          "401": { "description": "Invalid credentials" }
        }
      }
    },
    "/users": {
      "get": {
        "tags": ["Users"],
        "summary": "List all users (with pagination)",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "parameters": [
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 10 } },
          { "name": "offset", "in": "query", "schema": { "type": "integer", "default": 0 } }
        ],
        "responses": {
          "200": {
            "description": "List of users",
            "content": { "application/json": { "schema": { "type": "array", "items": { "$ref": "#/components/schemas/User" } } } }
          }
        }
      },
      "post": {
        "tags": ["Users"],
        "summary": "Create a new user",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateUserRequest" }
            }
          }
        },
        "responses": {
          "201": { "description": "User created" },
          "400": { "description": "Validation error" }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get user by ID",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": {
            "description": "User found",
            "content": { "application/json": { "schema": { "$ref": "#/components/schemas/User" } } }
          },
          "404": { "description": "User not found" }
        }
      },
      "put": {
        "tags": ["Users"],
        "summary": "Update user by ID",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateUserRequest" }
            }
          }
        },
        "responses": {
          "200": { "description": "User updated" },
          "404": { "description": "User not found" }
        }
      },
      "delete": {
        "tags": ["Users"],
        "summary": "Soft-delete user by ID (admin only)",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": { "description": "User deleted" },
          "403": { "description": "Forbidden - admin only" },
          "404": { "description": "User not found" }
        }
      }
    },
    "/users/with-audit": {
      "post": {
        "tags": ["Users"],
        "summary": "Create user with audit log (transaction)",
        "security": [{ "ApiKeyAuth": [] }, { "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateUserRequest" }
            }
          }
        },
        "responses": {
          "201": { "description": "User created with audit log" },
          "400": { "description": "Validation error" }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "ApiKeyAuth": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-KEY"
      },
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": { "type": "integer", "example": 1 },
          "name": { "type": "string", "example": "John Doe" },
          "email": { "type": "string", "example": "john@example.com" },
          "age": { "type": "integer", "example": 30 },
          "role": { "type": "string", "example": "user" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "CreateUserRequest": {
        "type": "object",
        "required": ["name", "email"],
        "properties": {
          "name": { "type": "string", "example": "Jane Doe" },
          "email": { "type": "string", "example": "jane@example.com" },
          "age": { "type": "integer", "example": 25 },
          "password": { "type": "string", "example": "secret123" },
          "role": { "type": "string", "example": "user", "enum": ["user", "admin"] }
        }
      },
      "UpdateUserRequest": {
        "type": "object",
        "properties": {
          "name": { "type": "string", "example": "Jane Updated" },
          "email": { "type": "string", "example": "jane.updated@example.com" },
          "age": { "type": "integer", "example": 26 }
        }
      }
    }
  }
}`
