package models

import "time"

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Verificated bool      `json:"-"`
	VerifToken  string    `json:"-"`
	Expired     time.Time `json:"-"`
}
