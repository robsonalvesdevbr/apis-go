package entity

import (
	"errors"

	"github.com/robsonalvesdevbr/apis-go/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password must not be empty")
	}

	uuidNew, errGuid := entity.NewID()
	if errGuid != nil {
		return nil, errGuid
	}
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return nil, errHash
	}
	return &User{
		ID:       uuidNew,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
