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
		AllowedOrigins: []string{"http://*", "https://*"},
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
	bookingsService := service.NewBookingService(bookingsRepo, userRepo, timeslotRepo)
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
			r.Get("/{id}", timeslotHandler.GetTimeslotByID)
			r.Get("/by-court", timeslotHandler.GetTimeslotsByCourt)
			r.Get("/", timeslotHandler.GetAllTimeslots)
			r.Group(func(r chi.Router) {
				r.Use(userAuthHandler.Authenticate)
				r.Use(userAuthHandler.RequireAdmin)
				r.Put("/{id}", timeslotHandler.UpdateTimeslot)
				r.Post("/", timeslotHandler.CreateTimeslot)
				r.Delete("/{id}", timeslotHandler.DeleteTimeslot)
			})
		})

		// Bookings
		r.Route("/bookings", func(r chi.Router) {
			r.Get("/{id}", bookingHandler.GetBookingByID)
			r.Get("/", bookingHandler.GetAllBookings)
			r.Group(func(r chi.Router) {
				r.Use(userAuthHandler.Authenticate)
				r.Post("/", bookingHandler.CreateBooking)
				r.Put("/{id}", bookingHandler.UpdateBooking)
				r.Delete("/{id}", bookingHandler.DeleteBooking)
			})
		})

		// Location routes
		r.Route("/location", func(r chi.Router) {
			r.Get("/", locationHandler.GetAllLocations)
			r.Group(func(r chi.Router) {
				r.Use(userAuthHandler.Authenticate)
				r.Use(userAuthHandler.RequireAdmin)
				r.Post("/", locationHandler.CreateLocation)
				r.Put("/", locationHandler.UpdateLocation)
				r.Delete("/{id}", locationHandler.DeleteLocation)
				r.Post("/{id}/image", locationHandler.UploadLocationImage)
			})
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser) // Keep this public for user registration
			r.Group(func(r chi.Router) {
				r.Use(userAuthHandler.Authenticate)
				r.Get("/{id}", userHandler.GetUserByID)
				r.Group(func(r chi.Router) {
					r.Use(userAuthHandler.RequireAdmin)
					r.Get("/", userHandler.GetAllUsers)
					r.Put("/", userHandler.UpdateUser)
					r.Delete("/{id}", userHandler.DeleteUser)
				})
			})
		})

		// Court routes
		r.Route("/courts", func(r chi.Router) {
			r.Get("/by-location", courtHandler.GetCourtByLocation)
			r.Get("/{id}", courtHandler.GetCourtByID)
			r.Get("/", courtHandler.GetAllCourts)
			r.Group(func(r chi.Router) {
				r.Use(userAuthHandler.Authenticate)
				r.Use(userAuthHandler.RequireAdmin)
				r.Post("/", courtHandler.CreateCourt)
				r.Put("/", courtHandler.UpdateCourt)
				r.Delete("/{id}", courtHandler.DeleteCourt)
			})
		})

		// Authentication routes
		r.Post("/login", userAuthHandler.Login)
		r.Post("/refresh", userAuthHandler.Refresh)
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"), // Matches the expected URL
	))
	r.Handle("/docs/*", http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))))

	return r
}
