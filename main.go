package main

import (
	swagger "github.com/gofiber/swagger"
	server "github.com/minacio00/easyCourt/Server"
	"github.com/minacio00/easyCourt/database"
)

func main() {
	// app := fiber.New()
	app := server.SetRoutes()
	app.Get("/swagger/*", swagger.HandlerDefault)

	database.Connectdb()
	app.Listen(":8080")
}
