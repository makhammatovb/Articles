package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/makhammatovb/Articles/internal/routes"
	"github.com/makhammatovb/Articles/internal/app"
)

// main function to start the server
func main() {
	// defines a command-line flag for the port number, defaulting to 8080 if not provided
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()
	// creates a new instance of the application and checks for errors
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	// sets up the routes using the chi router and the application instance
	r := routes.SetupRoutes(app)
	// configures and starts the HTTP server with specified timeouts
	server := &http.Server{
		Addr:       fmt.Sprintf(":%d", port),
		Handler:    r,
		IdleTimeout: time.Minute,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	// logs any fatal errors
	app.Logger.Printf("Listening on port %d", port)
	// starts the server and logs any errors
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}

