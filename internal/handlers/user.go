package handlers

import (
	"fiber-auth-api/internal/helper"
	"fiber-auth-api/internal/models"
	"fiber-auth-api/internal/repositories"
	"fiber-auth-api/internal/validation"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	app     models.Application
	dbModel *models.DbModel
}

func NewUserHandler(app models.Application, dbModel *models.DbModel) *UserHandler {
	return &UserHandler{app: app, dbModel: dbModel}
}

func (userHandler UserHandler) SignUpHandler(c fiber.Ctx) error {
	user := new(models.UserAuthenticationModel)

	if err := validation.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
		"username": true,
		"first_name": true,
		"last_name": true,
	}, user); err != nil {
		if invalidFieldErr, ok := validation.IsInvalidFieldError(err); ok {
			userHandler.app.SlogLogger.Error("Invalid field error", "error", invalidFieldErr)
			return userHandler.BadRequestFieldResponseError(c, fiber.Map{
				"email":    "user@example.com",
				"password": "yourPassword123!",
				"username": "johndoe",
				"first_name": "John",
				"last_name": "Doe",
			}, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}

		userHandler.app.SlogLogger.Error("Invalid json body", "error", err)
		return userHandler.BadRequestResponseError(c, fiber.Map{
			"email":    "user@example.com",
			"password": "yourPassword123!",
			"username": "johndoe",
			"first_name": "John",
			"last_name": "Doe",
		})
	}

	// v := validation.NewErrorValidator()
	// v.Check(user.Email != "", "email", "email must be provided")

	// if !v.IsValid() {
	// 	return userHandler.ValidationResponseError(c, fiber.Map{
	// 		"email":    "user@example.com",
	// 		"password": "SecurePass123!",
	// 	}, v.ValidationErrorField)
	// }

	ok, err := userHandler.dbModel.UserDbModel.IsUserExists(user.Email, user.Username)

	if err != nil {
		userHandler.app.SlogLogger.Error("Failed to check if user exists", "error", err)
		return userHandler.InternalServerErrorResponseError(c)
	}

	if ok {
		return userHandler.ConflictResponseError(c, "User already exists")
	}

	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		userHandler.app.SlogLogger.Error("Failed to hash password", "error", err)
	}

	userResponse := &repositories.UserCreateResponseModel{
		Email:     user.Email,
		PasswordHash:  hashedPassword,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  true,
	}

	err = userHandler.dbModel.UserDbModel.CreateUser(userResponse)

	if err != nil {
		userHandler.app.SlogLogger.Error("Something went wrong creating user", "error", err)
		return userHandler.InternalServerErrorResponseError(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "successful",
        "message": "User created successfully",
        "data":    userResponse.Email,
	})

}

func (userHandler UserHandler) SignInHandler(c fiber.Ctx) error {

	user := new(models.UserAuthenticationModel)

	if err := validation.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
	}, user); err != nil {
		if invalidFieldErr, ok := validation.IsInvalidFieldError(err); ok {
			return userHandler.BadRequestFieldResponseError(c, fiber.Map{
				"email":    "user@example.com",
				"password": "yourPassword123!",
			}, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}

		return userHandler.BadRequestResponseError(c, fiber.Map{
			"email":    "user@example.com",
				"password": "yourPassword123!",
				"username": "johndoe",
				"first_name": "John",
				"last_name": "Doe",
		})
	}

	if err := c.Bind().JSON(user); err != nil {
		userHandler.app.SlogLogger.Error("Failed to bind request body", "error", err)
		return userHandler.InternalServerErrorResponseError(c)
	}

	v := validation.NewErrorValidator()
	v.Check(user.Email != "", "email", "email must be provided")

	if !v.IsValid() {
		return userHandler.ValidationResponseError(c, fiber.Map{
			"email":    "user@example.com",
			"password": "SecurePass123!",
		}, v.ValidationErrorField)
	}

	return c.SendString(user.Email)
}

func (userHandler UserHandler) SignOutHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) ResetPasswordHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) GetAllUsersHandler(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (userHandler UserHandler) GetUserByIdHandler(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (userHandler UserHandler) GetUserByUsernameHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) GetUserByEmailHandler(c fiber.Ctx) error { return nil }
