package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/bhehar/bed-and-breakfast/internal/driver"
	"github.com/bhehar/bed-and-breakfast/internal/handlers"
	"github.com/bhehar/bed-and-breakfast/internal/helpers"
	"github.com/bhehar/bed-and-breakfast/internal/models"
	"github.com/bhehar/bed-and-breakfast/internal/render"
)

const (
	portNum = ":8080"
)

var (
	app     config.AppConfig
	session *scs.SessionManager
	infoLog *log.Logger
	errLog  *log.Logger
)

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close() 
	fmt.Printf("Listening on port %s\n", portNum)

	srv := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in prod
	app.InProd = false

	// setup loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrLog = infoLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProd

	app.Session = session

	//connect to database
	log.Println("Connecting to Database...")
	db, err := driver.ConnectSQL(`host=localhost port=5432 user=bali dbname=bookings`)
	if err != nil {
		log.Fatal("cannot connect to database! Drying.... :(")
	}
	log.Println("Connected to Database!")
	

	// create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("cannot create template cache. error: %v", err)
	}
	// assing to app config
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	// pass into render package
	render.NewRenderer(&app)
	return db, nil
}
