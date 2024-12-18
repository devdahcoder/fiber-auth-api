package handlers

import (
	"fiber-auth-api/internal/models"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	app     models.Application
	dbModel *models.DbModel
}

func NewUserHandler(app models.Application, dbModel *models.DbModel) *UserHandler {
	return &UserHandler{app: app, dbModel: dbModel}
}

func (userHandler UserHandler) SignInHandler(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (userHandler UserHandler) SignUpHandler(c fiber.Ctx) error {
	return nil
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

func (UserHandler UserHandler) GetUserByEmailHandler(c fiber.Ctx) error { return nil }
