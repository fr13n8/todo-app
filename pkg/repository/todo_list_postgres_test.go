package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fr13n8/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoListPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoListPostgres(db)

	type input struct {
		userId int
		list   todo.List
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
				userId: 1,
				list: todo.List{
					Id:          1,
					Title:       "Test title",
					Description: "Test description",
				},
			},
			wantId: 1,
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_lists").
					WithArgs(input.userId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty fields",
			input: input{
				userId: 1,
				list: todo.List{
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:    "Failed 2nd Insert",
			wantErr: true,
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_lists").
					WithArgs(input.userId, id).
					WillReturnError(errors.New("failed 2nd Insert"))

				mock.ExpectRollback()
			},
			input: input{
				userId: 1,
				list: todo.List{
					Title:       "title",
					Description: "description",
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.wantId)

			got, err := r.Create(testCase.input.userId, testCase.input.list)
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

func TestTodoListPostgres_GetById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoListPostgres(db)

	type input struct {
		listId int
		userId int
	}

	type mockBehavior func(input)

	testTable := []struct {
		name         string
		wantErr      bool
		input        input
		mockBehavior mockBehavior
		want         todo.List
	}{
		{
			name: "Ok",
			input: input{
				listId: 1,
				userId: 1,
			},
			want: todo.List{
				Id:          1,
				Title:       "title",
				Description: "description",
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(1, "title", "description")

				mock.ExpectQuery(`SELECT (.+) FROM todo_lists tl
									INNER JOIN users_lists ul ON (.+)
									WHERE (.+)`).
					WithArgs(input.userId, input.listId).
					WillReturnRows(rows)
			},
		},
		{
			name: "not found",
			input: input{
				listId: 1,
				userId: 1,
			},
			wantErr: true,
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery(`SELECT (.+) FROM todo_lists tl
									INNER JOIN users_lists ul ON (.+)
									WHERE (.+)`).
					WithArgs(input.userId, input.listId).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetById(testCase.input.listId, testCase.input.userId)
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

func TestTodoListPostgres_GetAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoListPostgres(db)

	type input struct {
		userId int
	}

	type mockBehavior func(input)

	testTable := []struct {
		name         string
		input        input
		mockBehavior mockBehavior
		wantErr      bool
		want         []todo.List
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
			},
			want: []todo.List{
				{
					Id:          1,
					Title:       "title",
					Description: "description",
				},
				{
					Id:          2,
					Title:       "title2",
					Description: "description2",
				},
				{
					Id:          3,
					Title:       "title3",
					Description: "description3",
				},
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow("1", "title", "description").
					AddRow("2", "title2", "description2").
					AddRow("3", "title3", "description3")

				mock.ExpectQuery(`SELECT (.+) FROM todo_lists tl 
										INNER JOIN users_lists ul ON (.+)
										WHERE (.+)`).
					WithArgs(input.userId).
					WillReturnRows(rows)
			},
		},
		{
			name: "No record found",
			input: input{
				userId: 1,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})
				mock.ExpectQuery(`SELECT (.+) FROM todo_lists tl 
										INNER JOIN users_lists ul ON (.+)
										WHERE (.+)`).
					WithArgs(input.userId).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetAll(testCase.input.userId)
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

func TestTodoListPostgres_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoListPostgres(db)

	type input struct {
		userId int
		listId int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		wantErr      bool
		input        input
		mockBehavior mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec(`DELETE FROM todo_lists tl 
									USING users_lists ul 
									WHERE (.+)`).
					WithArgs(input.userId, input.listId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "No record found",
			input: input{
				userId: 1,
				listId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec(`DELETE FROM todo_lists tl 
									USING users_lists ul 
									WHERE (.+)`).
					WithArgs(input.userId, input.listId).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			err := r.Delete(testCase.input.listId, testCase.input.userId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoListPostgres(db)

	type input struct {
		listId int
		userId int
		list   todo.UpdateListInput
	}

	type mockBehavior func(input)

	testTable := []struct {
		name         string
		wantErr      bool
		input        input
		mockBehavior mockBehavior
	}{
		{
			name: "Ok",
			input: input{
				userId: 1,
				listId: 1,
				list: todo.UpdateListInput{
					Title:       stringPointer("title"),
					Description: stringPointer("description"),
				},
			},
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs(input.list.Title, input.list.Description, input.listId, input.userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "OK_WithoutDescription",
			input: input{
				list: todo.UpdateListInput{
					Title: stringPointer("title"),
				},
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs(input.list.Title, input.listId, input.userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "OK_WithoutTitle",
			input: input{
				list: todo.UpdateListInput{
					Description: stringPointer("description"),
				},
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs(input.list.Description, input.listId, input.userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "OK_NoInputFields",
			input: input{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs(input.listId, input.userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			err := r.Update(testCase.input.listId, testCase.input.userId, testCase.input.list)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
