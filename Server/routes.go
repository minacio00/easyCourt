package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/minacio00/easyCourt/handlers"
)

func SetRoutes() *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			// Log the error
			log.Printf("Error: %v", err)
			log.Printf("Method: %v", c.Method())
			// You can also log more details if needed, such as request information
			// log.Printf("Error: %v, Path: %s, Method: %s", err, c.Path(), c.Method())
			return err
		}
		return nil
	})
	app.Use(cors.New(cors.ConfigDefault))
	// app.Use(func(c *fiber.Ctx) error {
	// 	c.Set("Access-Control-Allow-Origin", "*")
	// 	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// 	c.Set("Access-Control-Allow-Headers", "Content-Type")
	// 	return c.Next()
	// })

	app.Post("/tenant", handlers.CreateTenant)
	app.Get("/tenant/:id", handlers.ReadTenant)
	app.Get("/tenants/", handlers.GetTenants)
	app.Patch("/tentant/:id", handlers.UpdateTenant)
	app.Delete("/tenant/:id", handlers.DeleteTenant)

	app.Post("/quadra", handlers.CreateQuadra)
	app.Get("/quadra/:id", handlers.ReadQuadra)
	app.Get("/quadras/", handlers.GetQuadras)
	app.Patch("/quadra/:id", handlers.UpdateQuadra)
	app.Delete("/quadra/:id", handlers.DeleteQuadra)

	app.Post("/reserva", handlers.CreateReserva)
	app.Get("/reserva/:id", handlers.ReadReserva)
	app.Get("/reservas/", handlers.GetReservas)
	app.Patch("/reserva/:id", handlers.UpdateReserva)
	app.Delete("/reserva/:id", handlers.DeleteReserva)

	app.Post("/cliente", handlers.CreateCliente)
	app.Get("/cliente/:id", handlers.ReadCliente)
	app.Get("/clientes/", handlers.GetClientes)
	app.Patch("/cliente/:id", handlers.UpdateCliente)
	app.Delete("/cliente/:id", handlers.DeleteCliente)

	app.Post("/clube", handlers.CreateClube)
	app.Get("/clube/:id", handlers.ReadClube)
	app.Get("/clubes/", handlers.GetClubes)
	app.Patch("/clube/:id", handlers.UpdateClube)
	app.Delete("/clube/:id", handlers.DeleteClube)
	return app
}
