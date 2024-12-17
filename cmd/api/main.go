package main

import (
	"database/sql"
	"fiber-auth-api/internal/database"
	"fiber-auth-api/internal/logger"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

type Application struct {
	fiberApp *fiber.App
	log      *slog.Logger
	db       *sql.DB
}

func main() {

	logger.InitializeLogger(logger.SlogLogConfig{
		Level: slog.LevelDebug,
		JSON:  false,
	})

	log := logger.GetLogger()

	db := database.GetPsqlDatabase()
	defer func() {
		if err := db.ClosePsqlDb(); err != nil {
			log.Error(fmt.Sprintf("Error closing database: %v", err))
		}
	}()

	app := Application{
		fiberApp: route(),
		log:      log,
		db:       db.GetPsqlDB(),
	}

	logger.Error(fmt.Sprintf("Error closing database: %v", app.fiberApp.Listen(":3000")))

}
