package store

import (
	"fmt"
	"testmod/internal/app/model"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	model.TodoStore
	model.UserStore
	model.RecoveryStore
}

func NewStore(DbURL string) (*Store, error) {
	db, err := sqlx.Open("postgres", DbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	return &Store{
		TodoStore: NewTodoStore(db),
		UserStore: NewUserStore(db),
		RecoveryStore: NewRecoveryStore(db),
	}, nil
}
