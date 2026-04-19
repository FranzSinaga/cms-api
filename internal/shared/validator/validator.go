package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("password", validatePassword)
}

// Validate validates a struct based on validation tags
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// GetValidationErrors converts validator errors to a readable format
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())
			errors[field] = getErrorMessage(e)
		}
	}

	return errors
}

// getErrorMessage returns a user-friendly error message for each validation error
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param())
	case "password":
		return "Password must be at least 8 characters and contain uppercase, lowercase, and number"
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// validatePassword is a custom validator for password strength
// Password must be at least 8 characters and contain at least one uppercase, one lowercase, and one number
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum length 8
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Check for at least one digit
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}
