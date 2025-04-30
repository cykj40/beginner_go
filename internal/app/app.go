package app

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cykj40/beginner_go/internal/api"
	"github.com/cykj40/beginner_go/internal/store"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(db, migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             db,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
