package rest

import "platform-exercise/internal/entities"

type UserResponse struct {
	Data  *entities.User `json:"data,omitempty"`
	Error *string        `json:"error,omitempty"`
}

type LogInUserResponse struct {
	Data  map[string]string `json:"data,omitempty"`
	Error *string           `json:"error,omitempty"`
}
