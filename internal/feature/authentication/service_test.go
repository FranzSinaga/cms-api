package authentication

import (
	"errors"
	"testing"
	"time"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of the Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtConfig := config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 24 * time.Hour,
	}
	service := NewAuthService(mockRepo, jwtConfig)

	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1234",
		Name:     "John Doe",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)

	err := service.Register(req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Verify that CreateUser was called with correct data
	calls := mockRepo.Calls
	assert.Equal(t, 1, len(calls))
	user := calls[0].Arguments.Get(0).(*domain.User)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, "admin", user.Role)
	assert.NotEqual(t, req.Password, user.Password) // Password should be hashed
}

func TestRegister_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtConfig := config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 24 * time.Hour,
	}
	service := NewAuthService(mockRepo, jwtConfig)

	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1234",
		Name:     "John Doe",
	}

	expectedErr := errors.New("database error")
	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(expectedErr)

	err := service.Register(req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtConfig := config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 24 * time.Hour,
	}
	service := NewAuthService(mockRepo, jwtConfig)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test1234"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "John Doe",
		Role:     "admin",
	}

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "Test1234",
	}

	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	token, userResponse, err := service.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, userResponse)
	assert.Equal(t, existingUser.Email, userResponse.Email)
	assert.Equal(t, existingUser.Name, userResponse.Name)
	assert.Equal(t, existingUser.Role, userResponse.Role)

	// Verify JWT token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, existingUser.ID, claims["user_id"])
	assert.Equal(t, existingUser.Email, claims["email"])
	assert.Equal(t, existingUser.Role, claims["role"])

	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtConfig := config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 24 * time.Hour,
	}
	service := NewAuthService(mockRepo, jwtConfig)

	req := &domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "Test1234",
	}

	mockRepo.On("FindByEmail", req.Email).Return(nil, errors.New("user not found"))

	token, userResponse, err := service.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, userResponse)
	assert.Contains(t, err.Error(), "invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtConfig := config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 24 * time.Hour,
	}
	service := NewAuthService(mockRepo, jwtConfig)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("CorrectPassword123"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "John Doe",
		Role:     "admin",
	}

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword123",
	}

	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	token, userResponse, err := service.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, userResponse)
	assert.Contains(t, err.Error(), "invalid email or password")
	mockRepo.AssertExpectations(t)
}
