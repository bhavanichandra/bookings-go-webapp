package main

import (
	// Standard Library imports
	"fmt"
	"log"
	"net/http"
	"time"

	// Third Party Library Imports
	"github.com/alexedwards/scs/v2"

	// Application Packages
	"github.com/bhavanichandra/bookings/pkg/config"
	"github.com/bhavanichandra/bookings/pkg/handlers"
	"github.com/bhavanichandra/bookings/pkg/render"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

//main is the main application function
func main() {

	//Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	render.NewTemplates(&app)
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	fmt.Println(fmt.Sprintf("Listening to port %s", PORT))
	serve := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
