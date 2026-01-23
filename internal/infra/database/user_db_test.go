package database

import (
	"testing"

	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUser_DB_CreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "johndoe@example.com", "securepassword")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.NoError(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Password, userFound.Password)
}
