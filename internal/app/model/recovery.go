package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AccountRecovery struct {
	Id         string    `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Hash       string    `json:"hash" db:"ver_hash"`
	Expiration time.Time `json:"expiration" db:"expiration"`
}

type RecoveryStore interface {
	Create(a *AccountRecovery) (string, error)
	Get(id string) (string, string, error)
}

func (a *AccountRecovery) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}
