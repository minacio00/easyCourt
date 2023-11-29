package main

import (
	"log"

	swagger "github.com/gofiber/swagger"
	server "github.com/minacio00/easyCourt/Server"
	"github.com/minacio00/easyCourt/database"
)

func main() {
	// app := fiber.New()
	// app.Use(cors.New())

	app := server.SetRoutes()
	app.Get("/swagger/*", swagger.HandlerDefault)

	database.Connectdb()
	log.Fatal(app.Listen(":8080"))

}
