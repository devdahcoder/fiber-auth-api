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
	UserDbModel *repositories.UserRepository
}

func NewDbModel(userRepository *repositories.UserRepository) *DbModel {
	return &DbModel{UserDbModel: userRepository}
}

func (dbModel DbModel) GetUserRepository() *repositories.UserRepository {
	return dbModel.UserDbModel
}