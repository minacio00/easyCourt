package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Message string `json:"message"`
}

func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "Hello, World!"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	http.ListenAndServe(":8080", r)
}
