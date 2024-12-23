package validation

import (
	"encoding/json"
	"errors"
	"fiber-auth-api/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v3"
)

type ValidationErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	ValidationErrorField []ValidationErrorField
}

func NewErrorValidator() *ValidationError {
	return &ValidationError{ValidationErrorField: make([]ValidationErrorField, 0)}
}

func (validatorError *ValidationError) AddError(field string, message string) {
	validatorError.ValidationErrorField = append(validatorError.ValidationErrorField, ValidationErrorField{Field: field, Message: message})
}

func (validatorError *ValidationError) Check(ok bool, field string, message string) {
	if !ok {
		validatorError.AddError(field, message)
	}
}

func (validatorError *ValidationError) IsValid() bool {
	return len(validatorError.ValidationErrorField) == 0
}

type InvalidFieldError struct {
	Fields []string
}

func NewInvalidFieldError(field []string) *InvalidFieldError {
	return &InvalidFieldError{Fields: field}
}

func IsInvalidFieldError(err error) (*InvalidFieldError, bool) {
	var invalidFieldErr *InvalidFieldError
	ok := errors.As(err, &invalidFieldErr)
	return invalidFieldErr, ok
}

func InvalidFieldValidation(c fiber.Ctx, expectedFields map[string]bool) error {

	body := c.BodyRaw()
	var rawFields map[string]interface{}

	if err := json.Unmarshal(body, &rawFields); err != nil {
		return err
	}

	var unknownFields []string
	for field := range rawFields {
		if _, exists := expectedFields[field]; !exists {
			unknownFields = append(unknownFields, field)
		}
	}

	if len(unknownFields) > 0 {
		return NewInvalidFieldError(unknownFields)
	}

	return json.Unmarshal(body, &models.UserAuthenticationModel{})

}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("unknown field(s): %v", e.Fields)
}
