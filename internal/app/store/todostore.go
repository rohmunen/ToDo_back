package store

import (
	"fmt"
	"math/rand"
	"testmod/internal/app/model"

	"github.com/jmoiron/sqlx"
)

type TodoStore struct {
	*sqlx.DB
}

func NewTodoStore(db *sqlx.DB) *TodoStore {
	return &TodoStore{
		DB: db,
	}
}

func (s *TodoStore) Todo(id string) (*model.Todo, error) {
	var t model.Todo
	if err := s.Get(&t, `SELECT * FROM todos WHERE id = $1`, id); err != nil {
		fmt.Println("err")
		return &model.Todo{}, fmt.Errorf("error getting todo: %w", err)
	}
	return &t, nil
}

func (s *TodoStore) Todos(id string) ([]model.Todo, error) {
	var t []model.Todo
	if err := s.Select(&t, `SELECT * FROM todos WHERE user_id = $1`, id); err != nil {
		fmt.Println("err")
		return []model.Todo{}, fmt.Errorf("error getting todos: %w", err)
	}
	return t, nil
}

func (s *TodoStore) Create(t *model.Todo) (*model.Todo, error) {
	if err := s.QueryRowx(`INSERT INTO todos (title, body, user_id, isdone, isfavourite) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		t.Title,
		t.Body,
		t.UserId,
		t.IsDone,
		t.IsFavourite).Scan(&t.Id); err != nil {
		return &model.Todo{}, fmt.Errorf("error creating todo: %w", err)
	}
	return t, nil
}

func (s *TodoStore) Delete(id string) error {
	if _, err := s.Exec(`DELETE FROM todos_public WHERE todo_id = $1`, id); err != nil {
		return fmt.Errorf("error deleting todo: %w", err)
	}
	if _, err := s.Exec(`DELETE FROM todos WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting todo: %w", err)
	}
	return nil
}

func (s *TodoStore) Update(t *model.Todo) error {
	if err := s.QueryRowx(`UPDATE todos SET title = $1, body = $2, isdone = $3, isfavourite = $4 WHERE id = $5 RETURNING *`,
		t.Title,
		t.Body,
		t.IsDone,
		t.IsFavourite,
		t.Id).StructScan(t); err != nil {
		return fmt.Errorf("error updating todo: %w", err)
	}
	return nil
}

func (s *TodoStore) TodoPublic(id string) (string, error) {
	alphaNumRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	link := make([]rune, 60)
	for i := 0; i < 60; i++ {
		link[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	linkString := string(link)
	t, err := s.Todo(id)
	if err != nil {
		return "", fmt.Errorf("error getting todos: %w", err)
	}
	if err := s.QueryRowx(`INSERT INTO todos_public (todo_id, link_string) VALUES ($1, $2) RETURNING id`,
		id,
		linkString).Scan(&t.Id); err != nil {
		return "", fmt.Errorf("error creating todo: %w", err)
	}
	return linkString, nil
}

func (s *TodoStore) TodoPublicGet(link string) (*model.Todo, error) {
	var id string
	if err := s.QueryRowx(`SELECT todo_id FROM todos_public WHERE link_string = $1`, link).Scan(&id); err != nil {
		return &model.Todo{}, fmt.Errorf("error getting todo_id: %w", err)
	}
	t, err := s.Todo(id)
	if err != nil {
		return &model.Todo{}, fmt.Errorf("error getting todo: %w", err)
	}
	return t, nil
}
