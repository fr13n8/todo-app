package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fr13n8/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
		item   todo.Item
	}
	type mockBehavior func(input input, id int)

	testTable := []struct {
		name         string
		input        input
		mockBehavior mockBehavior
		wantId       int
		wantErr      bool
	}{
		// {
		// 	name: "OK",
		// 	input: input{
		// 		listId: 1,
		// 		item: todo.Item{
		// 			Title:       "Test title",
		// 			Description: "Test description",
		// 		},
		// 	},
		// 	wantId: 2,
		// 	mockBehavior: func(input input, id int) {
		// 		mock.ExpectBegin()

		// 		rows := sqlmock.NewRows([]string{"wantId"}).AddRow(id)
		// 		mock.ExpectQuery("INSERT INTO todo_items").
		// 			WithArgs(input.item.Title, input.item.Description).
		// 			WillReturnRows(rows)

		// 		mock.ExpectExec("INSERT INTO lists_items").
		// 			WithArgs(input.listId, id).
		// 			WillReturnResult(sqlmock.NewResult(1, 1))

		// 		mock.ExpectCommit()
		// 	},
		// },
		{
			name: "Empty Fields",
			input: input{
				listId: 1,
				item: todo.Item{
					Title:       "",
					Description: "description",
				},
			},
			mockBehavior: func(input input, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(input.item.Title, input.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		// {
		// 	name: "Failed 2nd Insert",
		// 	input: input{
		// 		listId: 1,
		// 		item: todo.Item{
		// 			Title:       "title",
		// 			Description: "description",
		// 		},
		// 	},
		// 	wantId: 2,
		// 	mockBehavior: func(input input, id int) {
		// 		mock.ExpectBegin()

		// 		rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
		// 		mock.ExpectQuery("INSERT INTO todo_items").
		// 			WithArgs(input.item.Title, input.item.Description).WillReturnRows(rows)

		// 		mock.ExpectExec("INSERT INTO lists_items").WithArgs(input.listId, id).
		// 			WillReturnError(errors.New("insert error"))

		// 		mock.ExpectRollback()
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.wantId)

			got, err := r.Create(testCase.input.listId, testCase.input.item)
			if testCase.wantErr {
				logrus.Println(got, err)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantId, got)
			}
		})
	}
}
