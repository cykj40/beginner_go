package store

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	Logger *log.Logger
}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{
		Logger: log.New(log.Writer(), "[WorkoutStore] ", log.Ldate|log.Ltime),
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Println("Handling GET workout by ID")
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "this is the workout id %d\n", workoutID)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Println("Handling POST workout")
	fmt.Fprintf(w, "created a workout\n")
}
