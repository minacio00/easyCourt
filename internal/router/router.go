// easyCourt/internal/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/minacio00/easyCourt/docs"
	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/handler"
	"github.com/minacio00/easyCourt/internal/repository"
	"github.com/minacio00/easyCourt/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userAuthHandler := handler.NewUserAuthHandler(userService)

	// Define routes
	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Get("/", userHandler.GetAllUsers)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Authentication route
		r.Post("/login", userAuthHandler.Login)
	})

	r.Post("/login", userAuthHandler.Login)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"), // Matches the expected URL
	))
	r.Handle("/docs/*", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	return r
}
