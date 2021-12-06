package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int    `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password,omitempty" db:"password"`
	EncryptedPassword string `json:"-" db:"encrypted_password"`
}

type UserStore interface {
	Create(user *User) error
	FindById(userId string) (*User, error)
	FindByEmail(email string) (*User, error)
	UpdatePassword(userId string, password string) (*User, error)
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (u *User) BeforeCreate() error {
	fmt.Println("in before create: ", u.Password)
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}
