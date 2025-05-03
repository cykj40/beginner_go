package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cykj40/beginner_go/internal/utils"

	"github.com/cykj40/beginner_go/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to get workout", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid workout update id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	// Read the entire request body
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		wh.logger.Printf("ERROR: decodingUpdateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Log the raw request for debugging
	jsonBody, _ := json.Marshal(requestBody)
	wh.logger.Printf("Received update request: %s", string(jsonBody))

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
		wh.logger.Printf("ERROR: updateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update the workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
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
		wh.logger.Printf("ERROR: strconv.ParseInt: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid workout delete id"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete workout"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
