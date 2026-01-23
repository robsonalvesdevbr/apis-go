package entity

import (
	"errors"

	"github.com/robsonalvesdevbr/apis-go/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserName     = errors.New("user name must not be empty")
	ErrInvalidUserEmail    = errors.New("user email must not be empty")
	ErrInvalidUserPassword = errors.New("user password must not be empty")
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {
	uuidNew, errGuid := entity.NewID()
	if errGuid != nil {
		return nil, errGuid
	}

	user := &User{
		ID:       uuidNew,
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return nil, errHash
	}

	user.Password = string(hashedPassword)

	return user, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	if u.Email == "" {
		return ErrInvalidUserEmail
	}
	if u.Password == "" {
		return ErrInvalidUserPassword
	}
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
