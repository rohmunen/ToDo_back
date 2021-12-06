package store

import (
	"testmod/internal/app/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RecoveryStore struct {
	*sqlx.DB
}

func NewRecoveryStore(db *sqlx.DB) *RecoveryStore {
	return &RecoveryStore{
		DB: db,
	}
}

func (s *RecoveryStore) Create(verificationRow *model.VerificationRow) (string, error) {
	if err := verificationRow.Validate(); err != nil {
		return "", err
	}
	if err := s.QueryRowx(`INSERT INTO email_ver_hash (email, ver_hash, expiration) VALUES ($1, $2, $3) RETURNING id`,
		verificationRow.Email,
		verificationRow.Hash,
		verificationRow.Expiration,
	).Scan(
		&verificationRow.Id,
	); err != nil {
		return "", err
	}
	return verificationRow.Id, nil
}

func (s *RecoveryStore) Get(verificationId string) (string, string, error) {
	var hash string
	var email string
	if err := s.QueryRowx(`SELECT ver_hash, email FROM email_ver_hash WHERE id = $1`,
		verificationId,
	).Scan(&hash, &email); err != nil {
		return "", "", err
	}
	return hash, email, nil
}
