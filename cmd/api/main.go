package main

import (
	"fiber-auth-api/internal/database"
	"github.com/gofiber/fiber/v3"
	"log/slog"
)

type dbConfig struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	fiberApp *fiber.App
	log      *slog.Logger
	dbConfig *dbConfig
}

func main() {

	db := database.InitializeDb()

	defer db.Close()

	//dbConfig := dbConfig{
	//	port: 8080,
	//	env:  "development",
	//	db: struct{ dsn string }{
	//		dsn: "user:password@tcp(localhost:3306)/your_database",
	//	},
	//}

	app := application{
		fiberApp: route(),
		//dbConfig: &dbConfig,
	}

	err = app.fiberApp.Listen(":3000")

	if err != nil {
		return
	}
}
