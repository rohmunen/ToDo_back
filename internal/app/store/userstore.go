package store

import (
	"fmt"
	"testmod/internal/app/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStore struct {
	*sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

func (s *UserStore) Create(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if err := user.BeforeCreate(); err != nil {
		return err
	}
	if err := s.QueryRowx(`INSERT INTO users (email, encrypted_password, name) VALUES ($1, $2, $3) RETURNING id`,
		user.Email,
		user.EncryptedPassword,
		user.Name,
	).Scan(&user.Id); err != nil {
		return fmt.Errorf("error inserting new user into db: %w", err)
	}
	return nil
}

func (s *UserStore) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := s.DB.QueryRowx(`SELECT id, email, encrypted_password FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return user, nil
}

func (s *UserStore) FindById(id string) (*model.User, error) {
	user := &model.User{}
	if err := s.DB.QueryRowx(`SELECT id, email, encrypted_password FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return user, nil
}

func (s *UserStore) UpdatePassword(id string, password string) (*model.User, error) {
	user := &model.User{}
	if err := s.DB.QueryRowx(`SELECT id, email, encrypted_password FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, fmt.Errorf("error finding user by id: %w", err)
	}
	user.Password = password
	user.BeforeCreate()
	user.Sanitize()
	if err := s.DB.QueryRowx(`UPDATE users SET encrypted_password=$1 WHERE id = $2`,
		user.EncryptedPassword,
		user.Id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}
	return user, nil
}
