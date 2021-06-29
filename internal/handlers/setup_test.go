package handlers

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/models"
	"github.com/sunil206b/smart_booking/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var appConfig config.AppConfig
var session *scs.SessionManager
var templatesPath = "../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	//what am I going to put in the session
	gob.Register(models.Reservation{})

	//Change this to true when in the production
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatalf("Error while creating template cache: %v\n", err)
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = true
	render.NewTemplates(&appConfig)

	rhHandler := NewRouteHandler(&appConfig)
	NewHandler(rhHandler)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	//router.Use(NoSurf)
	router.Use(SessionLoad)

	router.Get("/", Handler.Home)
	router.Get("/about", Handler.About)
	router.Get("/contact", Handler.Contact)

	router.Get("/generals-quarters", Handler.Generals)
	router.Get("/majors-suite", Handler.Majors)

	router.Get("/search-availability", Handler.Availability)
	router.Post("/search-availability", Handler.PostAvailability)
	router.Post("/search-availability-json", Handler.AvailabilityJSON)

	router.Get("/make-reservations", Handler.Reservation)
	router.Post("/make-reservations", Handler.PostReservation)
	router.Get("/reservation-summary", Handler.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache function creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatesPath))
	if err != nil {
		return nil, errors.New("Error while looking for pages " + err.Error())
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, errors.New("Error while generating template set " + err.Error())
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
		if err != nil {
			return nil, errors.New("Error while checking for layout file " + err.Error())
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatesPath))
			if err != nil {
				return nil, errors.New("Error while parsing layout file " + err.Error())
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
