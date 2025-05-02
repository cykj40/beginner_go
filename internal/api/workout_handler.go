package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cykj40/beginner_go/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
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

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to get workout", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	fmt.Printf("Updating workout with ID: %s\n", paramsWorkoutID)

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing workout ID: %v\n", err)
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		fmt.Printf("Error fetching workout for update: %v\n", err)
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		fmt.Printf("Workout with ID %d not found\n", workoutID)
		http.NotFound(w, r)
		return
	}

	// Read the entire request body
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Printf("Error decoding update request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the raw request for debugging
	jsonBody, _ := json.Marshal(requestBody)
	fmt.Printf("Received update request: %s\n", string(jsonBody))

	if title, ok := requestBody["title"].(string); ok {
		existingWorkout.Title = title
	}

	if desc, ok := requestBody["description"].(string); ok {
		existingWorkout.Description = desc
	}

	if duration, ok := requestBody["duration_minutes"].(float64); ok {
		existingWorkout.DurationMinutes = int(duration)
	}

	if calories, ok := requestBody["calories_burned"].(float64); ok {
		existingWorkout.CaloriesBurned = int(calories)
	}

	// Only update entries if provided in the request body
	if entries, ok := requestBody["entries"].([]interface{}); ok && len(entries) > 0 {
		// Clear existing entries and add new ones
		existingWorkout.Entries = []store.WorkoutEntry{}

		for _, entryData := range entries {
			if entryMap, ok := entryData.(map[string]interface{}); ok {
				entry := store.WorkoutEntry{}

				if name, ok := entryMap["exercise_name"].(string); ok {
					entry.ExerciseName = name
				}

				if sets, ok := entryMap["sets"].(float64); ok {
					entry.Sets = int(sets)
				}

				if reps, ok := entryMap["reps"].(float64); ok {
					repsInt := int(reps)
					entry.Reps = &repsInt
				}

				if duration, ok := entryMap["duration_seconds"].(float64); ok {
					durationInt := int(duration)
					entry.DurationSeconds = &durationInt
				}

				if weight, ok := entryMap["weight"].(float64); ok {
					entry.Weight = &weight
				}

				if notes, ok := entryMap["notes"].(string); ok {
					entry.Notes = &notes
				}

				if index, ok := entryMap["order_index"].(float64); ok {
					entry.OrderIndex = int(index)
				}

				existingWorkout.Entries = append(existingWorkout.Entries, entry)
			}
		}
	} else {
		fmt.Println("No entries provided in update, keeping existing entries")
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		fmt.Printf("Error updating workout in store: %v\n", err)
		http.Error(w, "failed to update the workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}

// HandleUpdateWorkoutByID handles updates to a workout by its ID
// It delegates to HandleUpdateWorkout for the actual update logic
func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	wh.HandleUpdateWorkout(w, r)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
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

	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
