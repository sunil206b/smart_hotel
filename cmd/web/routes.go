package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//router := mux.NewRouter()
	//router.HandleFunc("/", handlers.Handler.Home).Methods(http.MethodGet)
	//router.HandleFunc("/about", handlers.Handler.About).Methods(http.MethodGet)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(NoSurf)
	router.Use(SessionLoad)

	router.Get("/", handlers.Handler.Home)
	router.Get("/about", handlers.Handler.About)
	router.Get("/contact", handlers.Handler.Contact)

	router.Get("/generals-quarters", handlers.Handler.Generals)
	router.Get("/majors-suite", handlers.Handler.Majors)

	router.Get("/search-availability", handlers.Handler.Availability)
	router.Post("/search-availability", handlers.Handler.PostAvailability)
	router.Post("/search-availability-json", handlers.Handler.AvailabilityJSON)

	router.Get("/make-reservations", handlers.Handler.Reservation)
	router.Post("/make-reservations", handlers.Handler.PostReservation)
	router.Get("/reservation-summary", handlers.Handler.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}