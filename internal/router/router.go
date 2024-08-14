// easyCourt/internal/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/handler"
	"github.com/minacio00/easyCourt/internal/repository"
	"github.com/minacio00/easyCourt/internal/service"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userAuthHandler := handler.NewUserAuthHandler(userService)

	// Define routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUserByID)
		r.Get("/", userHandler.GetAllUsers)
		r.Put("/", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	r.Post("/login", userAuthHandler.Login)

	return r
}
