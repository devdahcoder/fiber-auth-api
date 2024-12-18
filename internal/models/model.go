package models

import (
	"database/sql"
	"fiber-auth-api/internal/repositories"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

type Application struct {
	FiberApp   *fiber.App
	SlogLogger *slog.Logger
	PsqlDb     *sql.DB
}

func (app *Application) NewApplication(fiber *fiber.App, slogLogger *slog.Logger,
	psqlDb *sql.DB) *Application {
	return &Application{
		FiberApp:   fiber,
		SlogLogger: slogLogger,
		PsqlDb:     psqlDb,
	}
}

type DbModel struct {
	userDbModel *repositories.UserRepository
}

func NewDbModel(user *repositories.UserRepository) *DbModel {
	return &DbModel{userDbModel: user}
}
