package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minacio00/easyCourt/internal/tenant"
)

type Response struct {
	Message string `json:"message"`
}

func StartServer() {
	tenant.Migrate()

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(log.Writer(), "", log.LstdFlags), NoColor: false}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "Hello, World!"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Post("/", tenant.CreateTenantHandler)

	log.Println("starting server on :8008")
	http.ListenAndServe(":8080", r)
}
