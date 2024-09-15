// easyCourt/internal/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/minacio00/easyCourt/docs"
	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/handler"
	"github.com/minacio00/easyCourt/internal/repository"
	"github.com/minacio00/easyCourt/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	r.Use(middleware.Logger)

	// Set up repositories
	userRepo := repository.NewUserRepository(db.GetDB())
	locationRepo := repository.NewLocationRepository(db.GetDB())
	courtRepo := repository.NewCourtRepository(db.GetDB())
	bookingsRepo := repository.NewBookingRepository(db.GetDB())
	timeslotRepo := repository.NewTimeslotRepository(db.GetDB())

	// Set up services
	userService := service.NewUserService(userRepo)
	locationService := service.NewLocationService(locationRepo)
	courtService := service.NewCourtService(courtRepo)
	bookingsService := service.NewBookingService(bookingsRepo)
	timeslotService := service.NewTimeslotService(timeslotRepo, courtRepo)

	// Set up handlers
	userHandler := handler.NewUserHandler(userService)
	userAuthHandler := handler.NewUserAuthHandler(userService)
	locationHandler := handler.NewLocationHandler(locationService)
	courtHandler := handler.NewCourtHandler(courtService)
	bookingHandler := handler.NewBookingHandler(bookingsService)
	timeslotHandler := handler.NewTimeslotHandler(timeslotService)

	// Define routes
	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/timeslots", func(r chi.Router) {
			r.Post("/", timeslotHandler.CreateTimeslot)
			r.Get("/", timeslotHandler.GetAllTimeslots)
			r.Put("/", timeslotHandler.UpdateTimeslot)
			r.Get("/{id}", timeslotHandler.GetTimeslotByID)
			r.Delete("/{id}", timeslotHandler.DeleteTimeslot)
		})

		r.Route("/bookings", func(r chi.Router) {
			r.Post("/", bookingHandler.CreateBooking)
			r.Get("/", bookingHandler.GetAllBookings)
			r.Put("/", bookingHandler.UpdateBooking)
			r.Delete("/", bookingHandler.DeleteBooking)
		})

		// Location routes
		r.Route("/location", func(r chi.Router) {
			r.Post("/", locationHandler.CreateLocation)
			r.Get("/", locationHandler.GetAllLocations)
			r.Put("/", locationHandler.UpdateLocation)
			r.Delete("/{id}", locationHandler.DeleteLocation)
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Get("/", userHandler.GetAllUsers)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Court routes
		r.Route("/courts", func(r chi.Router) {
			r.Post("/", courtHandler.CreateCourt)
			r.Get("/", courtHandler.GetAllCourts)
			r.Get("/{id}", courtHandler.GetCourtByID)
			r.Put("/", courtHandler.UpdateCourt)
			r.Delete("/{id}", courtHandler.DeleteCourt)
		})

		// Authentication route
		r.Post("/login", userAuthHandler.Login)
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"), // Matches the expected URL
	))
	r.Handle("/docs/*", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	return r
}
