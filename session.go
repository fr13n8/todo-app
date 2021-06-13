package todo

import "github.com/google/uuid"

type Session struct {
	Id           int       `json:"id" db:"id"`
	UserId       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	UUID         uuid.UUID `json:"uuid" db:"uuid"`
}
