package entities

import "time"

type User struct {
	ID       *string    `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Password string     `json:"-"`
	Birthday *time.Time `json:"birthday"`
}
