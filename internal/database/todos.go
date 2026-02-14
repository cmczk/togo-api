package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
		id, title, description, completed,
		created_at, updated_at
	FROM todos
	ORDER BY created_at DESC;`

	if err := s.db.Select(&todos, query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return todos, nil
}

func (s *TodoStore) GetByID(id int) (*models.Todo, error) {
	const op = "database.todos.GetByID"

	var todo models.Todo

	query := `
	SELECT
		id, title, description, completed,
		created_at, updated_at
	FROM todos
	WHERE id = $1;`

	if err := s.db.Get(&todo, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo not found")
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &todo, nil
}

func (s *TodoStore) Create(data models.CreateTodoInput) (*models.Todo, error) {
	const op = "database.todos.Create"

	var todo models.Todo

	query := `
	INSERT INTO
		todos (title, description, completed)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, completed, created_at, updated_at;`

	if err := s.db.QueryRowx(
		query, data.Title, data.Description, data.Completed,
	).StructScan(&todo); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &todo, nil
}

func (s *TodoStore) Update(id int, data models.UpdateTodoInput) (*models.Todo, error) {
	const op = "database.todos.Update"

	todo, err := s.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if data.Title != nil {
		todo.Title = *data.Title
	}

	if data.Description != nil {
		todo.Description = *data.Description
	}

	if data.Completed != nil {
		todo.Completed = *data.Completed
	}

	query := `
	UPDATE todos
	SET title = $1, description = $2, completed = $3, updated_at = $4
	WHERE id = $5
	RETURNING id, title, description, completed, created_at, updated_at;`

	var updatedTodo models.Todo

	if err := s.db.QueryRowx(
		query, todo.Title, todo.Description, todo.Completed, time.Now(), id,
	).StructScan(&updatedTodo); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedTodo, nil
}

func (s *TodoStore) Delete(id int) error {
	const op = "database.todos.Delete"

	query := "DELETE FROM todos WHERE id = $1;"

	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
