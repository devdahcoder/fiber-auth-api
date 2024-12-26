package handlers

import (
	"errors"
	"fiber-auth-api/internal/helper"
	"fiber-auth-api/internal/models"
	"fiber-auth-api/internal/repositories"
	"fiber-auth-api/internal/types"
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
	
var (
	exampleEmail = "user@example.com"
	requestExample = fiber.Map{
		"email":    exampleEmail,
		"password": exampleEmail,
		"username": "johndoe",
		"first_name": "John",
		"last_name": "Doe",
	}
)

func (userHandler UserHandler) SignUpHandler(c fiber.Ctx) error {
	
	user := new(repositories.UserCreateModel)
	userHandler.ValidateSignUp(c, user)

	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		userHandler.app.SlogLogger.Error("Failed to hash password", "error", err)
		return userHandler.InternalServerErrorResponseError(c)
	}

	userResponse := &repositories.UserCreateDbModel{
		Email:     user.Email,
		PasswordHash:  hashedPassword,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  true,
	}

	err = userHandler.dbModel.UserDbModel.CreateUser(userResponse)

	if err != nil {
		if errors.Is(err, types.ErrDuplicateUser) {
			return userHandler.ConflictResponseError(c, "User already exists")
		}
		return userHandler.InternalServerErrorResponseError(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "successful",
        "message": "User created successfully",
        "data":    userResponse.Email,
	})

}

func (userHandler UserHandler) SignInHandler(c fiber.Ctx) error {

	user := new(repositories.UserCreateModel)

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
	users, err := userHandler.dbModel.UserDbModel.GetAllUsers()
	if err != nil {
		return userHandler.InternalServerErrorResponseError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "successful",
        "message": "All users fetched successfully",
        "data":    users,
	})
}

func (userHandler UserHandler) GetUserByIdHandler(c fiber.Ctx) error {

	userId := c.Params("id")

	user, err := userHandler.dbModel.UserDbModel.FindUserById(userId)
	if err != nil {
		return userHandler.NotFoundResponseError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "successful",
        "message": "All user fetched successfully",
        "data":    user,
	})
}

func (userHandler UserHandler) GetUserByUsernameHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) GetUserByEmailHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) ValidateSignUp(c fiber.Ctx, user *repositories.UserCreateModel) error {
	if err := validation.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
		"username": true,
		"first_name": true,
		"last_name": true,
	}, user); err != nil {
		if invalidFieldErr, ok := validation.IsInvalidFieldError(err); ok {
			userHandler.app.SlogLogger.Error("Invalid field error", "error", invalidFieldErr)
			return userHandler.BadRequestFieldResponseError(c, requestExample, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}

		userHandler.app.SlogLogger.Error("Invalid json body", "error", err)
		return userHandler.BadRequestResponseError(c, requestExample)
	}

	v := validation.NewErrorValidator()
	v.Check(user.Email != "", "email", "email must be provided")
	v.Check(user.Password != "", "password", "password must be provided")
	v.Check(user.Username != "", "username", "username must be provided")
	v.Check(user.FirstName != "", "first_name", "first name must be provided")
	v.Check(user.LastName != "", "last_name", "last name must be provided")
	
	if !v.IsValid() {
		return userHandler.ValidationResponseError(c, requestExample, v.ValidationErrorField)
	}

	return nil

}