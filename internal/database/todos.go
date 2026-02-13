package database

import (
	"fmt"

	"github.com/cmczk/todo-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type TodoStore struct {
	db *sqlx.DB
}

func NewTodoStore(db *sqlx.DB) *TodoStore {
	return &TodoStore{db: db}
}

func (s *TodoStore) GetAll() (todos []models.Todo, err error) {
	const op = "database.todos.GetAll"

	query := `
	SELECT
		id, title, description,
		created_at, updated_at
	FROM todos
	ORDER BY created_at DESC;`

	if err := s.db.Select(&todos, query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return todos, nil
}
