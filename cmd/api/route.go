package main

import "github.com/gofiber/fiber/v3"

func route() *fiber.App {

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app

}
