package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cykj40/beginner_go/internal/app"
	"github.com/cykj40/beginner_go/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to start the server on")
	flag.Parse()

	log.Println("Starting application...")
	app, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	log.Println("Setting up routes...")
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server starting on port %d", port)
	log.Printf("Available endpoints:")
	log.Printf("  POST /users/register")
	log.Printf("  POST /users/login")
	log.Printf("  GET  /health")
	log.Printf("  GET  /workouts/{id}")
	log.Printf("  POST /workouts")
	log.Printf("  PUT  /workouts/{id}")
	log.Printf("  DELETE /workouts/{id}")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
