package configs

import (
	"errors"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) LoadConfig(path string) (*conf, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*conf), args.Error(1)
}

func TestLoadConfigSuccess(t *testing.T) {
	mockConfig := new(MockConfig)

	mockConfig.On("LoadConfig", "/tmp/test").Return(&conf{
		DBDriver:      "mysql",
		DBHost:        "localhost",
		DBPort:        "3306",
		DBUser:        "root",
		DBPassword:    "root",
		DBName:        "testdb",
		WebServerPort: "8080",
		JwtSecret:     "my-secret-key",
		JwtExpiresIn:  3600,
		LogLevel:      "info",
		AuthToken:     jwtauth.New("HS256", []byte("my-secret-key"), nil),
	}, nil)
	config, err := mockConfig.LoadConfig("/tmp/test")

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "mysql", config.DBDriver)
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "3306", config.DBPort)
	assert.Equal(t, "root", config.DBUser)
	assert.Equal(t, "root", config.DBPassword)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "8080", config.WebServerPort)
	assert.Equal(t, "my-secret-key", config.JwtSecret)
	assert.Equal(t, 3600, config.JwtExpiresIn)
	assert.Equal(t, "info", config.LogLevel)
	assert.NotNil(t, config.AuthToken)

	mockConfig.AssertExpectations(t)
}

func TestLoadConfigFileNotFound(t *testing.T) {
	mockConfig := new(MockConfig)

	mockConfig.On("LoadConfig", "/tmp/test").Return(nil, errors.New("config file not found"))
	config, err := mockConfig.LoadConfig("/tmp/test")

	assert.Error(t, err)
	assert.Nil(t, config)

	mockConfig.AssertExpectations(t)
}
