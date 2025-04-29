package api

import (
	"net/http"

	"github.com/cykj40/beginner_go/internal/store"
)

type WorkoutHandler struct {
	store *store.WorkoutHandler
}

func NewWorkoutHandler() *WorkoutHandler {
	storeHandler := store.NewWorkoutHandler()
	return &WorkoutHandler{
		store: storeHandler,
	}
}

func (h *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	h.store.HandleGetWorkoutByID(w, r)
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	h.store.HandleCreateWorkout(w, r)
}

// Add handler methods here as needed
