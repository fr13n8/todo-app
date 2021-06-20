package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fr13n8/todo-app/structs"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoItemPostgres(db)

	type input struct {
		listId int
		item   structs.Item
	}
	type mockBehavior func(input input, id int)

	testTable := []struct {
		name         string
		input        input
		mockBehavior mockBehavior
		wantId       int
		wantErr      bool
	}{
		{
			name: "OK",
			input: input{
				listId: 1,
				item: structs.Item{
					Id:          1,
					Title:       "Test title",
					Description: "Test description",
				},
			},
			wantId: 1,
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"wantId"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(input.item.Title, input.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(input.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: input{
				listId: 1,
				item: structs.Item{
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(input.item.Title, input.item.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed 2nd Insert",
			input: input{
				listId: 1,
				item: structs.Item{
					Id:          1,
					Title:       "title",
					Description: "description",
				},
			},
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(input.item.Title, input.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(input.listId, id).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.wantId)

			got, err := r.Create(testCase.input.listId, testCase.input.item)
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

func TestTodoItemPostgres_GetById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoItemPostgres(db)

	type input struct {
		itemId int
		userId int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		mockBehavior mockBehavior
		want         structs.Item
		wantErr      bool
	}{
		{
			name: "OK",
			input: input{
				userId: 1,
				itemId: 1,
			},
			want: structs.Item{
				Id:          1,
				Title:       "title",
				Description: "description",
				Done:        false,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title", "description", false)

				mock.ExpectQuery(`SELECT (.+) FROM todo_items ti
									INNER JOIN lists_items li on (.+)
									INNER JOIN users_lists ul on (.+)
									WHERE (.+)`).
					WithArgs(input.itemId, input.userId).
					WillReturnRows(rows)
			},
		},
		{
			name: "not found",
			input: input{
				userId: 1,
				itemId: 1,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery(`SELECT (.+) FROM todo_items ti 
									INNER JOIN lists_items li on (.+)
									INNER JOIN users_lists ul on (.+)
									WHERE (.+)`).
					WithArgs(input.itemId, input.userId).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetById(testCase.input.userId, testCase.input.itemId)
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

func TestTodoItemPostgres_GetAll(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoItemPostgres(db)

	type input struct {
		listId int
		userId int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		wantErr      bool
		mockBehavior mockBehavior
		want         []structs.Item
	}{
		{
			name: "OK",
			input: input{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow("1", "title", "description", false).
					AddRow("2", "title1", "description1", false).
					AddRow("3", "title2", "description2", true)

				mock.ExpectQuery(`SELECT (.+) FROM todo_items ti 
									INNER JOIN lists_items li on (.+)
									INNER JOIN users_lists ul on (.+)
									WHERE (.+)`).
					WithArgs(input.listId, input.userId).
					WillReturnRows(rows)
			},
			want: []structs.Item{
				{
					Id:          1,
					Title:       "title",
					Description: "description",
					Done:        false,
				},
				{
					Id:          2,
					Title:       "title1",
					Description: "description1",
					Done:        false,
				},
				{
					Id:          3,
					Title:       "title2",
					Description: "description2",
					Done:        true,
				},
			},
		},
		{
			name: "no records",
			input: input{
				listId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery(`SELECT (.+) FROM todo_items ti 
									INNER JOIN lists_items li on (.+)
									INNER JOIN users_lists ul on (.+)
									WHERE (.+)`).
					WithArgs(input.listId, input.userId).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err := r.GetAll(testCase.input.listId, testCase.input.userId)
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

func TestTodoItemPostgres_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoItemPostgres(db)

	type input struct {
		userId int
		itemId int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		wantErr      bool
		mockBehavior mockBehavior
		input        input
	}{
		{
			name: "Ok",
			mockBehavior: func(input input) {
				mock.ExpectExec(`DELETE FROM todo_items ti 
									USING lists_items li, users_lists ul 
									WHERE (.+)`).
					WithArgs(input.userId, input.itemId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: input{
				userId: 1,
				itemId: 1,
			},
		},
		{
			name: "No record found",
			mockBehavior: func(input input) {
				mock.ExpectExec(`DELETE FROM todo_items ti 
									USING lists_items li, users_lists ul 
									WHERE (.+)`).
					WithArgs(input.userId, input.itemId).
					WillReturnError(sql.ErrNoRows)
			},
			input: input{
				userId: 1,
				itemId: 1,
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			err := r.Delete(testCase.input.userId, testCase.input.itemId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "smock")

	r := NewTodoItemPostgres(db)

	type input struct {
		item   structs.UpdateItemInput
		userId int
		itemId int
	}

	type mockBehavior func(input input)

	testTable := []struct {
		name         string
		input        input
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Ok",
			input: input{
				item: structs.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
				itemId: 1,
				userId: 1,
			},
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(input.item.Title, input.item.Description, input.item.Done, input.itemId, input.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "OK_WithoutDone",
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(input.item.Title, input.item.Description, input.itemId, input.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemId: 1,
				userId: 1,
				item: structs.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDoneAndDescription",
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(input.item.Title, input.itemId, input.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemId: 1,
				userId: 1,
				item: structs.UpdateItemInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mockBehavior: func(input input) {
				mock.ExpectExec("UPDATE todo_items ti SET FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(input.itemId, input.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemId: 1,
				userId: 1,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			err := r.Update(testCase.input.userId, testCase.input.itemId, testCase.input.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
