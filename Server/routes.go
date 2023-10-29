package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/handlers"
)

func SetRoutes() *fiber.App {
	app := fiber.New()

	app.Post("/tenant", handlers.CreateTenant)
	app.Get("/tenant/:id", handlers.ReadTenant)
	app.Get("/tenants/", handlers.GetTenants)
	app.Patch("/tentant/:id", handlers.UpdateTenant)
	app.Delete("/tenant/:id", handlers.DeleteTenant)

	return app
}
