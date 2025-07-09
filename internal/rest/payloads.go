package rest

import "time"

type CreateUserPayload struct {
	// user name
	Name string `json:"name" binding:"required,min=1"`
	// user email
	Email string `json:"email" binding:"required,email"`
	// user password
	Password string `json:"password" binding:"required,min=8"`
	// user birthday
	Birthday *time.Time `json:"birthday" time_format:"2006-01-02"`
} // @name CreateUserPayload

type LogInUserPayload struct {
	// user email
	Email string `json:"email" binding:"required,email"`
	// user password
	Password string `json:"password" binding:"required,min=8"`
} // @name LogInUserPayload

type UserPatchPayload struct {
	// user name
	Name string `json:"name" binding:"required,min=1"`
	// user birthday
	Birthday *time.Time `json:"birthday" time_format:"2006-01-02"`
} // @name UserPatchPayload
