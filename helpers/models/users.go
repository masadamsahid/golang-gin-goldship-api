package models

import (
	common "github.com/masadamsahid/golang-gin-goldship-api/helpers/commons"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"string"`
	common.BaseEntity
}

type Profile struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	User        User   `json:"user"`
	PhoneNumber string `json:"phone_number"`
	common.BaseEntity
}
