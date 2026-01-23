package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser_Success(t *testing.T) {
	user, err := NewUser("John Doe", "john.doe@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@example.com", user.Email)
}

func TestUser_ValidatePassword_IsValid(t *testing.T) {
	user, err := NewUser("Jane Doe", "jane.doe@example.com", "password123")
	assert.NoError(t, err)

	isValid := user.ValidatePassword("password123")
	assert.True(t, isValid)
}

func TestUser_NewUser_Password_IsInvalid(t *testing.T) {
	user, err := NewUser("John Doe", "john.doe@example.com", "password123")
	assert.NoError(t, err)

	isInvalid := user.ValidatePassword("wrongpassword")
	assert.False(t, isInvalid)
}

// validar se nao é hash
func TestUser_NewUser_Password_IsHashed(t *testing.T) {
	user, err := NewUser("John Doe", "john.doe@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEqual(t, "password123", user.Password)
}
