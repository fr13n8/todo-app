package todo

import "github.com/google/uuid"

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" binding:"required" db:"name"`
	UserName string `json:"username" binding:"required" db:"username"`
	Password string `json:"password" binding:"required" db:"password"`
}

type SignInInput struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required" db:"name"`
	UserName string `json:"username" binding:"required" db:"username"`
	Password string `json:"password" binding:"required" db:"password"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type Session struct {
	Id           int       `json:"id" db:"id"`
	UserId       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	UUID         uuid.UUID `json:"uuid" db:"uuid"`
}
