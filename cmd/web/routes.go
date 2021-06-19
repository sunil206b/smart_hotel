package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sunil206b/smart_booking/pkg/config"
	"github.com/sunil206b/smart_booking/pkg/handlers"
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

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}