package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/handlers"
	"github.com/sunil206b/smart_booking/internal/models"
	"github.com/sunil206b/smart_booking/internal/render"
	"log"
	"net/http"
	"time"
)

var (
	appConfig config.AppConfig
	session   *scs.SessionManager
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Failed to run with error %v\n", err)
	}
	//http.HandleFunc("/", handlers.Handler.Home)
	//http.HandleFunc("/about", handlers.Handler.About)
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatalln(err)
	//}

	srv := &http.Server{
		Handler:      routes(&appConfig),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("server will be running on port %v\n", 8080)
	log.Fatalln(srv.ListenAndServe())
}

func run() error {
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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("Error while creating template cache: %v\n", err)
		return err
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false
	render.NewTemplates(&appConfig)

	rhHandler := handlers.NewRouteHandler(&appConfig)
	handlers.NewHandler(rhHandler)
	return nil
}
