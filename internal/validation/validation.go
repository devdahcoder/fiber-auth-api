package validation

import (
	"encoding/json"
	"errors"
	"fiber-auth-api/internal/types"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Json field validation
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

// Json value validation
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

func InvalidFieldValidation(c fiber.Ctx, expectedFields map[string]bool, dataModel interface{}) error {
    body := c.BodyRaw()
    var rawFields map[string]interface{}

    if err := json.Unmarshal(body, &rawFields); err != nil {
        return err
    }

    unknownFields := findUnknownFields(rawFields, expectedFields)
    if len(unknownFields) > 0 {
        return NewInvalidFieldError(unknownFields)
    }

    if err := json.Unmarshal(body, &dataModel); err != nil {
        return err
    }

    return nil
}

func findUnknownFields(rawFields map[string]interface{}, expectedFields map[string]bool) []string {
    var unknownFields []string
    for field := range rawFields {
        if _, exists := expectedFields[field]; !exists {
            unknownFields = append(unknownFields, field)
        }
    }
    return unknownFields
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("unknown field(s): %v", e.Fields)
}

type QueryValidationError struct {
    Parameter string `json:"parameter"`
    Value     string `json:"value"`
    Message   string `json:"message"`
}

type QueryValidator struct {
    paramPatterns map[string]*regexp.Regexp
    typeValidators map[string]types.QueryTypeValidationFunction
}

func NewQueryValidator() *QueryValidator {
    qv := &QueryValidator{
        paramPatterns: make(map[string]*regexp.Regexp),
        typeValidators: make(map[string]types.QueryTypeValidationFunction),
    }
    
    qv.AddParamPattern("default", `^[a-zA-Z][a-zA-Z0-9_]*$`)
    
    qv.typeValidators["number"] = func(v string) bool {
        matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, v)
        return matched
    }

    qv.typeValidators["boolean"] = func(v string) bool {
        v = strings.ToLower(v)
        return v == "true" || v == "false" || v == "1" || v == "0"
    }

    qv.typeValidators["date"] = func(v string) bool {
        matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, v)
        return matched
    }
    
    return qv
}

func (qv *QueryValidator) AddParamPattern(name, pattern string) error {
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return fmt.Errorf("invalid pattern for %s: %v", name, err)
    }
    qv.paramPatterns[name] = regex
    return nil
}

func (qv *QueryValidator) AddTypeValidator(name string, validator types.QueryTypeValidationFunction) {
    qv.typeValidators[name] = validator
}

func (qv *QueryValidator) ValidateQuery(c fiber.Ctx, rules map[string]string) []QueryValidationError {
    var errors []QueryValidationError
    
    queries := c.Queries()
    
    for param, value := range queries {
        if !qv.validateParamName(param) {
            errors = append(errors, QueryValidationError{
                Parameter: param,
                Value:     value,
                Message:   "invalid parameter name format",
            })
            continue
        }
        
        expectedType, exists := rules[param]
        if !exists {
            errors = append(errors, QueryValidationError{
                Parameter: param,
                Value:     value,
                Message:   "unexpected parameter",
            })
            continue
        }
        
        if !qv.validateParamValue(value, expectedType) {
            errors = append(errors, QueryValidationError{
                Parameter: param,
                Value:     value,
                Message:   fmt.Sprintf("invalid value for type %s", expectedType),
            })
        }
    }
    
    return errors
}

func (qv *QueryValidator) validateParamName(param string) bool {
    pattern, exists := qv.paramPatterns["default"]
    if !exists {
        return true 
    }
    return pattern.MatchString(param)
}

func (qv *QueryValidator) validateParamValue(value, expectedType string) bool {
    validator, exists := qv.typeValidators[expectedType]
    if !exists {
        return true 
    }
    return validator(value)
}