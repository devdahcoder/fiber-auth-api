package main

import (
	"fiber-auth-api/internal/database"
	"fiber-auth-api/internal/logger"
	"fiber-auth-api/internal/models"
	"fiber-auth-api/internal/route"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

func main() {

	fiberApp := fiber.New(fiber.Config{})

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

	app := models.Application{
		FiberApp:   fiberApp,
		SlogLogger: log,
		PsqlDb:     db.GetPsqlDB(),
	}

	route.SetupRoutes(app)

	logger.Error(fmt.Sprintf("Error starting up application: %v", fiberApp.Listen(":3000")))

}
