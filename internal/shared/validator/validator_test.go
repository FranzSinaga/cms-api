package validator

import (
	"testing"

	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegisterRequest_Success(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1234",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.NoError(t, err)
}

func TestValidateRegisterRequest_InvalidEmail(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "invalid-email",
		Password: "Test1234",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "email")
}

func TestValidateRegisterRequest_MissingEmail(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "",
		Password: "Test1234",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "email")
}

func TestValidateRegisterRequest_PasswordTooShort(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
}

func TestValidateRegisterRequest_PasswordMissingUppercase(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "test1234",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
}

func TestValidateRegisterRequest_PasswordMissingLowercase(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "TEST1234",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
}

func TestValidateRegisterRequest_PasswordMissingNumber(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "TestTest",
		Name:     "John Doe",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
}

func TestValidateRegisterRequest_NameTooShort(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1234",
		Name:     "J",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "name")
}

func TestValidateRegisterRequest_MissingName(t *testing.T) {
	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "Test1234",
		Name:     "",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "name")
}

func TestValidateLoginRequest_Success(t *testing.T) {
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "anypassword",
	}

	err := Validate(req)
	assert.NoError(t, err)
}

func TestValidateLoginRequest_InvalidEmail(t *testing.T) {
	req := &domain.LoginRequest{
		Email:    "invalid-email",
		Password: "anypassword",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "email")
}

func TestValidateLoginRequest_MissingPassword(t *testing.T) {
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	err := Validate(req)
	assert.Error(t, err)

	errors := GetValidationErrors(err)
	assert.Contains(t, errors, "password")
}
