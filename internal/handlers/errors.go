package handlers

import (
	"fiber-auth-api/internal/validation"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

func (userHandler UserHandler) ErrorResponse(c fiber.Ctx, code int, err error) error {
	userHandler.app.SlogLogger.Error(err.Error())
	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": err.Error(),
		"status":  http.StatusText(code),
	})
}

func (userHandler UserHandler) InternalServerErrorResponseError(c fiber.Ctx) error {
	err := fmt.Errorf("the server encountered a problem and could not process your request")
	return userHandler.ErrorResponse(c, fiber.StatusInternalServerError, err)
}

func (userHandler UserHandler) NotFoundResponseError(c fiber.Ctx) error {
	err := fmt.Errorf("resource could not be found")
	return userHandler.ErrorResponse(c, fiber.StatusNotFound, err)
}

func (userHandler UserHandler) UnauthorizedResponseError(c fiber.Ctx) error {
	err := fmt.Errorf("unauthorized access")
	return userHandler.ErrorResponse(c, fiber.StatusUnauthorized, err)
}

func (userHandler UserHandler) BadRequestResponseError(c fiber.Ctx, structure map[string]any) error {
	err := fmt.Errorf("invalid request format")

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message":         err.Error(),
		"details":         "Please provide a valid JSON request with the required fields",
		"status":          http.StatusText(fiber.StatusBadRequest),
		"code":            fiber.StatusBadRequest,
		"request_example": structure,
	})
}

func (userHandler UserHandler) BadRequestFieldResponseError(c fiber.Ctx, structure map[string]any, details map[string]any) error {
	err := fmt.Errorf("invalid field names in request")

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message":         err.Error(),
		"status":          http.StatusText(fiber.StatusBadRequest),
		"code":            fiber.StatusBadRequest,
		"details":         details,
		"request_example": structure,
	})
}

func (userHandler UserHandler) ValidationResponseError(c fiber.Ctx, structure map[string]any, errors []validation.ValidationErrorField) error {
	err := fmt.Errorf("validation failed")

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":            fiber.StatusBadRequest,
		"message":         err.Error(),
		"status":          http.StatusText(fiber.StatusBadRequest),
		"errors":          errors,
		"request_example": structure,
	})
}

func (userHandler UserHandler) ConflictResponseError(c fiber.Ctx, message string) error {
	err := fmt.Errorf("%s", message)
	return userHandler.ErrorResponse(c, fiber.StatusConflict, err)
}
