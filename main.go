package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
)

func main() {
	app := fiber.New()
	database.Connectdb()
	app.Listen(":8080")
}
