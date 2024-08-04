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

	tenantRepo := repository.NewTenantRepository(db.DB)
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)
	tenantAuthHandler := handler.NewTenantAuthHandler(tenantService)

	r.Route("/tenants", func(r chi.Router) {
		r.Post("/", tenantHandler.CreateTenant)
		r.Get("/", tenantHandler.GetAllTenants)
		r.Get("/{id}", tenantHandler.GetTenantByID)
		r.Put("/{id}", tenantHandler.UpdateTenant)
		r.Delete("/{id}", tenantHandler.DeleteTenant)
	})

	r.Post("/login", tenantAuthHandler.Login)

	return r
}
