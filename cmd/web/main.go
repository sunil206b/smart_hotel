package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/sunil206b/smart_booking/pkg/config"
	"github.com/sunil206b/smart_booking/pkg/handlers"
	"github.com/sunil206b/smart_booking/pkg/render"
	"log"
	"net/http"
	"time"
)

var (
	appConfig config.AppConfig
	session *scs.SessionManager
)

func main() {

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
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false
	render.NewTemplates(&appConfig)

	rhHandler := handlers.NewRouteHandler(&appConfig)
	handlers.NewHandler(rhHandler)

	//http.HandleFunc("/", handlers.Handler.Home)
	//http.HandleFunc("/about", handlers.Handler.About)
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatalln(err)
	//}

	srv := &http.Server{
		Handler: routes(&appConfig),
		Addr: ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
