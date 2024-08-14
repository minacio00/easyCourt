package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minacio00/easyCourt/config"
	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/router"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Easy Court API
// @version 1.0
// @description This is the api for easy court.
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()
	db.Init()

	r := router.NewRouter()

	chiRouter, ok := r.(*chi.Mux)
	if !ok {
		log.Fatalf("Could not type assert router to chi.Mux")
	}
	// Register swagger handler
	chiRouter.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", chiRouter); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
