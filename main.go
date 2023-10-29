package main

import (
	server "github.com/minacio00/easyCourt/Server"
	"github.com/minacio00/easyCourt/database"
)

func main() {
	// app := fiber.New()
	app := server.SetRoutes()
	database.Connectdb()
	app.Listen(":8080")
}
