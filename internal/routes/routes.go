package routes

import (
	"log"
	"net/http"

	"github.com/cykj40/beginner_go/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Add debug logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Health check
	r.Get("/health", app.HealthCheck)

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", app.UserHandler.HandleRegisterUser)
		r.Post("/login", app.TokenHandler.HandleCreateToken)
	})

	// Workout routes
	r.Route("/workouts", func(r chi.Router) {
		r.Get("/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
		r.Post("/", app.WorkoutHandler.HandleCreateWorkout)
		r.Put("/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	})

	// Print registered routes
	log.Println("Registered routes:")
	log.Println("GET    /health")
	log.Println("POST   /users/register")
	log.Println("POST   /users/login")
	log.Println("GET    /workouts/{id}")
	log.Println("POST   /workouts")
	log.Println("PUT    /workouts/{id}")
	log.Println("DELETE /workouts/{id}")

	return r
}
