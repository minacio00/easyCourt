// @title EasyCourt API
// @description This is a comprehensive API for the EasyCourt system.
// @version 1.0
// @BasePath /api/v1
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/config"
	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/router"
)

func main() {
	config.LoadConfig()
	db.Init()

	r := router.NewRouter()

	chiRouter, ok := r.(*chi.Mux)
	if !ok {
		log.Fatalf("Could not type assert router to chi.Mux")
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", chiRouter); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
