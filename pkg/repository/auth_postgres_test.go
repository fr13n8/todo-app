package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fr13n8/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoAuth_CreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewAuthPostgres(db)

	type input struct {
		user todo.SignUpInput
	}

	type mockBehavior func(input input, id int)

	testTable := []struct {
		name         string
		input        input
		wantErr      bool
		wantId       int
		mockBehavior mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				user: todo.SignUpInput{
					Name:     "name",
					UserName: "username",
					Password: "password",
				},
			},
			wantId: 1,
			mockBehavior: func(input input, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.user.Name, input.user.UserName, input.user.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty password field",
			input: input{
				user: todo.SignUpInput{
					Name:     "name",
					UserName: "username",
					Password: "",
				},
			},
			wantErr: true,
			mockBehavior: func(input input, id int) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.user.Name, input.user.UserName, input.user.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty name field",
			input: input{
				user: todo.SignUpInput{
					Name:     "",
					UserName: "username",
					Password: "password",
				},
			},
			wantErr: true,
			mockBehavior: func(input input, id int) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.user.Name, input.user.UserName, input.user.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "Empty username field",
			input: input{
				user: todo.SignUpInput{
					Name:     "name",
					UserName: "",
					Password: "password",
				},
			},
			wantErr: true,
			mockBehavior: func(input input, id int) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.user.Name, input.user.UserName, input.user.Password).
					WillReturnRows(rows)
			},
		},
		{
			name:    "Empty fields",
			wantErr: true,
			input: input{
				user: todo.SignUpInput{
					Name:     "",
					UserName: "",
					Password: "",
				},
			},
			mockBehavior: func(input input, id int) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(input.user.Name, input.user.UserName, input.user.Password).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.wantId)

			got, err := r.CreateUser(testCase.input.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoAuth_GetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewAuthPostgres(db)

	type input struct {
		username string
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		wantErr      bool
		mockBehavior mockBehavior
		want         todo.User
	}{
		{
			name: "Ok",
			input: input{
				username: "username",
			},
			want: todo.User{
				Id:       1,
				Name:     "name",
				UserName: "username",
				Password: "password",
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).
					AddRow(1, "name", "username", "password")

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(input.username).
					WillReturnRows(rows)
			},
		},
		{
			name: "No record",
			input: input{
				username: "username",
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"})

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(input.username).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetUser(testCase.input.username)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoAuth_CreateSession(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewAuthPostgres(db)

	type mockBehavior func(input todo.Session)

	testTable := []struct {
		name         string
		input        todo.Session
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "Ok",
			input: todo.Session{
				UserId:       1,
				RefreshToken: "refresh_token",
				UserAgent:    "user_agent",
				UUID:         "uuid",
			},
			mockBehavior: func(input todo.Session) {
				mock.ExpectExec("INSERT INTO users_sessions").
					WithArgs(input.UserId, input.UUID, input.RefreshToken, input.UserAgent).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Empty userId field",
			input: todo.Session{
				Id:           1,
				RefreshToken: "refresh_token",
				UserAgent:    "user_agent",
				UUID:         "uuid",
			},
			wantErr: true,
			mockBehavior: func(input todo.Session) {
				mock.ExpectExec("INSERT INTO users_sessions").
					WithArgs(input.UserId, input.UUID, input.RefreshToken, input.UserAgent).
					WillReturnError(errors.New("invalid session"))
			},
		},
		{
			name:    "Empty fields",
			wantErr: true,
			input:   todo.Session{},
			mockBehavior: func(input todo.Session) {
				mock.ExpectExec("INSERT INTO users_sessions").
					WithArgs(input.UserId, input.UUID, input.RefreshToken, input.UserAgent).
					WillReturnError(errors.New("empty session"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			err := r.CreateSession(testCase.input)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoAuth_GetSessionsByUserId(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewAuthPostgres(db)

	type input struct {
		id int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		wantErr      bool
		mockBehavior mockBehavior
		want         []todo.Session
	}{
		{
			name: "Ok",
			input: input{
				id: 1,
			},
			want: []todo.Session{
				{
					Id:           1,
					UserId:       1,
					RefreshToken: "refresh_token",
					UserAgent:    "user_agent",
					UUID:         "uuid",
				},
				{
					Id:           2,
					UserId:       1,
					RefreshToken: "refresh_token2",
					UserAgent:    "user_agent2",
					UUID:         "uuid2",
				},
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "user_agent", "uuid"}).
					AddRow(1, 1, "refresh_token", "user_agent", "uuid").
					AddRow(2, 1, "refresh_token2", "user_agent2", "uuid2")

				mock.ExpectQuery("SELECT (.+) FROM users_sessions WHERE (.+)").
					WithArgs(input.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "No record",
			input: input{
				id: 1,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "user_agent", "uuid"})

				mock.ExpectQuery("SELECT (.+) FROM users_sessions WHERE (.+)").
					WithArgs(input.id).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetSessionsByUserId(testCase.input.id)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoAuth_GetSessionByUUID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewAuthPostgres(db)

	type input struct {
		uuid string
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		wantErr      bool
		mockBehavior mockBehavior
		want         todo.Session
	}{
		{
			name: "Ok",
			input: input{
				uuid: "uuid",
			},
			want: todo.Session{
				Id:           1,
				UserId:       1,
				RefreshToken: "refresh_token",
				UserAgent:    "user_agent",
				UUID:         "uuid",
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "user_agent", "uuid"}).
					AddRow(1, 1, "refresh_token", "user_agent", "uuid")

				mock.ExpectQuery("SELECT (.+) FROM users_sessions WHERE (.+)").
					WithArgs(input.uuid).
					WillReturnRows(rows)
			},
		},
		{
			name: "No record",
			input: input{
				uuid: "uuid",
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "user_agent", "uuid"})

				mock.ExpectQuery("SELECT (.+) FROM users_sessions WHERE (.+)").
					WithArgs(input.uuid).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetSessionByUUID(testCase.input.uuid)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
