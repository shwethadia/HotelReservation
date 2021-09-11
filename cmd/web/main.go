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
	"github.com/shwethadia/HotelReservation/internal/driver"
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

	db, err := run()
	if err != nil {

		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener")

	listenForMail()

	/* 	from := "me@here.com"
	   	auth := smtp.PlainAuth("", from, "", "localhost")
	   	err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello World"))
	   	if err != nil {
	   		log.Println(err)
	   	}
	*/
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

func run() (*driver.DB, error) {

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	//Connect to Database
	log.Println("Connecting to database.....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=HotelReservation user=shwetha password=1234")
	if err != nil {

		log.Fatal("Cannot connect to database...")
	}

	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil

}
