package main

import (
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

	// app.Use(func(c *fiber.Ctx) error {
	// 	if err := c.Next(); err != nil {
	// 		// Log the error
	// 		log.Printf("Error: %v", err)
	// 		// You can also log more details if needed, such as request information
	// 		// log.Printf("Error: %v, Path: %s, Method: %s", err, c.Path(), c.Method())
	// 		return err
	// 	}
	// 	return nil
	// })

	app.Listen(":8080")

}
