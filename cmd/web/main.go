package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/handlers"
	"github.com/shwethadia/HotelReservation/internal/helpers"
	"github.com/shwethadia/HotelReservation/internal/models"
	"github.com/shwethadia/HotelReservation/internal/render"
)

const portNumber = ":8089"

var app config.AppConfig

var session *scs.SessionManager

var infoLog *log.Logger

var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {

		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{

		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	gob.Register(models.Reservation{})

	//Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return nil

}
