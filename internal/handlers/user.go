package handlers

import (
	"errors"
	"fiber-auth-api/internal/helper"
	"fiber-auth-api/internal/models"
	"fiber-auth-api/internal/repositories"
	"fiber-auth-api/internal/types"
	"fiber-auth-api/internal/validation"
	"time"

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
	signupRequestExample = fiber.Map{
		"email":    exampleEmail,
		"password": exampleEmail,
		"username": "johndoe",
		"first_name": "John",
		"last_name": "Doe",
	}
	signinRequestExample = fiber.Map{
		"email":    exampleEmail,
		"password": exampleEmail,
	}
)

type queryParams struct {
	value map[string]string
}

func (userHandler UserHandler) SignUpHandler(c fiber.Ctx) error {
	
	user := new(repositories.UserSignupModel)
	if err := validation.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
		"username": true,
		"first_name": true,
		"last_name": true,
	}, user); err != nil {
		if invalidFieldErr, ok := validation.IsInvalidFieldError(err); ok {
			userHandler.app.SlogLogger.Error("Invalid field error", "error", invalidFieldErr)
			return userHandler.BadRequestFieldResponseError(c, signupRequestExample, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}

		userHandler.app.SlogLogger.Error("Invalid json body", "error", err)
		return userHandler.BadRequestResponseError(c, signupRequestExample)
	}

	v := validation.NewErrorValidator()
	v.Check(user.Email != "", "email", "email must be provided")
	v.Check(user.Password != "", "password", "password must be provided")
	v.Check(user.Username != "", "username", "username must be provided")
	v.Check(user.FirstName != "", "first_name", "first name must be provided")
	v.Check(user.LastName != "", "last_name", "last name must be provided")
	
	if !v.IsValid() {
		return userHandler.ValidationResponseError(c, signupRequestExample, v.ValidationErrorField)
	}
	
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

	return userHandler.SuccessResponse(c, "User created successfully", user.Email)

}

func (userHandler UserHandler) SignInHandler(c fiber.Ctx) error {

	user := new(repositories.UserSigninModel)

	if err := validation.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
	}, user); err != nil {
		if invalidFieldErr, ok := validation.IsInvalidFieldError(err); ok {
			return userHandler.BadRequestFieldResponseError(c, signinRequestExample, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}
		userHandler.app.SlogLogger.Error("Invalid json body", "error", err)
		return userHandler.BadRequestResponseError(c, signinRequestExample)
	}

	v := validation.NewErrorValidator()
	v.Check(user.Email != "", "email", "email must be provided")
	v.Check(user.Password != "", "password", "password must be provided")

	if !v.IsValid() {
		return userHandler.ValidationResponseError(c, signinRequestExample, v.ValidationErrorField)
	}

	userResponse, err := userHandler.dbModel.UserDbModel.AuthenticateUser(user.Email)

	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return userHandler.NotFoundResponseError(c)
		}
		return userHandler.InternalServerErrorResponseError(c)
	}

	if err := helper.VerifyPassword(userResponse.PasswordHash, user.Password); err != nil {
		// "Invalid email or password"
		return userHandler.UnauthorizedResponseError(c)
	}

	token, err := helper.CreateToken(userResponse.Email)

	if err != nil {
		userHandler.app.SlogLogger.Error("Failed to create token", "error", err)
		return userHandler.InternalServerErrorResponseError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return userHandler.SuccessResponse(c, "User created successfully", token)

}

func (userHandler UserHandler) SignOutHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) ResetPasswordHandler(c fiber.Ctx) error { return nil }

func (userHandler UserHandler) GetAllUsersHandler(c fiber.Ctx) error {

	validator := validation.NewQueryValidator()
    
    // Define validation rules
    rules := map[string]string{
        "age":    "number",
        "status": "string",
        "search": "string",
    }
    
    // Validate query parameters
    if errors := validator.ValidateQuery(c, rules); len(errors) > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "errors": errors,
        })
    }

	users, metadata, err := userHandler.dbModel.UserDbModel.GetAllUsers()
	if err != nil {
		return userHandler.InternalServerErrorResponseError(c)
	}

	return userHandler.SuccessResponse(c, "All users fetched successfully", fiber.Map{
		"users":    users,
		"metadata": metadata,
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
