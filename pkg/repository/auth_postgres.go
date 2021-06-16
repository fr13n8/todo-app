package repository

import (
	"fmt"

	"github.com/fr13n8/todo-app"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.SignUpInput) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.UserName, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", usersTable)
	err := r.db.Get(&user, query, username)
	return user, err
}

func (r *AuthPostgres) CreateSession(input todo.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, uuid, refresh_token, uagent) VALUES ($1, $2, $3, $4)", usersSessionsTable)
	_, err := r.db.Exec(query, input.UserId, input.UUID, input.RefreshToken, input.UserAgent)
	return err
}

func (r *AuthPostgres) GetSessionsByUserId(userId int) ([]todo.Session, error) {
	var sessions []todo.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", usersSessionsTable)
	if err := r.db.Select(&sessions, query, userId); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *AuthPostgres) GetSessionByUUID(uuid string) (todo.Session, error) {
	var session todo.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE uuid=$1", usersSessionsTable)
	err := r.db.Get(&session, query, uuid)
	return session, err
}
