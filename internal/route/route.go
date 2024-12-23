package route

import (
	"fiber-auth-api/internal/handlers"
	"fiber-auth-api/internal/models"
	"fiber-auth-api/internal/repositories"
)

func SetupRoutes(app models.Application) {
	userRepository := repositories.NewUserRepository(app.PsqlDb)
	dbModel := models.NewDbModel(userRepository)
	userHandler := handlers.NewUserHandler(app, dbModel)

	apiV1 := app.FiberApp.Group("/api/v1")
	apiV1.Post("/signup", userHandler.SignUpHandler)
	apiV1.Post("/signin", userHandler.SignInHandler)
	apiV1.Post("/reset-password", userHandler.ResetPasswordHandler)
	apiV1.Get("/", userHandler.GetAllUsersHandler)
	apiV1.Get("/:id", userHandler.GetUserByIdHandler)
	apiV1.Get("/:username/", userHandler.GetUserByUsernameHandler)
	apiV1.Get("/:email/", userHandler.GetUserByEmailHandler)
}
