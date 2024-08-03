// easyCourt/internal/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define your routes here
	// r.Route("/your-route", func(r chi.Router) {
	//     r.Get("/", handler.YourHandler)
	// })

	return r
}
