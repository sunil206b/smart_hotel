package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/driver"
	"github.com/sunil206b/smart_booking/internal/handlers"
	"github.com/sunil206b/smart_booking/internal/helpers"
	"github.com/sunil206b/smart_booking/internal/models"
	"github.com/sunil206b/smart_booking/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	appConfig config.AppConfig
	session   *scs.SessionManager
	infoLog   *log.Logger
	errorLog  *log.Logger
)

func init() {
	gotenv.Load()
}

func main() {
	db, err := run()
	if err != nil {
		log.Fatalf("Failed to run with error %v\n", err)
	}
	defer db.SQL.Close()
	defer close(appConfig.MailChan)

	log.Println("Starting mail channel listener....")
	listenForMail()

	srv := &http.Server{
		Handler:      routes(&appConfig),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("server will be running on port %v\n", 8080)
	log.Fatalln(srv.ListenAndServe())
}

func run() (*driver.DB, error) {
	//what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan *models.MailData)
	appConfig.MailChan = mailChan

	//Change this to true when in the production
	appConfig.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	//connect to database
	log.Println("Connecting to database...")
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse Elephant SQL URL %v\n", err))
	}
	db, err := driver.ConnectPQSQL(pgURL)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect Elephant SQL %v\n", err))
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while creating template cache: %v\n", err))
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false
	render.NewRenderer(&appConfig)

	rhHandler := handlers.NewRouteHandler(&appConfig, db)
	handlers.NewHandler(rhHandler)

	helpers.NewHelpers(&appConfig)
	return db, nil
}
