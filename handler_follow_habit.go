package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/khangle880/living-plan/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFollowHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	habitIdStr := chi.URLParam(r, "habit_id")
	habitId, err := uuid.Parse(habitIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing habit id %v", err))
		return
	}

	type parameters struct {
		PromiseFrom    time.Time `json:"promise_from"`
		PromiseEnd     time.Time `json:"promise_end"`
		IncludedReport bool      `json:"included_report"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	followHabit, err := apiCfg.DB.CreateFollowHabit(r.Context(), database.CreateFollowHabitParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		PromiseFrom: params.PromiseFrom,
		PromiseEnd:  params.PromiseEnd,
		HabitID:     habitId,
		UserID:      user.ID,
		Processes:   []time.Time{},
	})

	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Couldn't follow habit %v", err))
		return
	}

	responseWithJSON(w, 201, followHabit)
}

// TODO: handle add and remove processes
func (apiCfg *apiConfig) handlerUpdateFollowHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	habitIdStr := chi.URLParam(r, "habit_id")
	habitId, err := uuid.Parse(habitIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing habit id %v", err))
		return
	}

	type parameters struct {
		PromiseFrom    time.Time `json:"promise_from"`
		PromiseEnd     time.Time `json:"promise_end"`
		IncludedReport *bool     `json:"included_report"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	includedReport := sql.NullBool{}
	if params.IncludedReport != nil {
		includedReport.Bool = *params.IncludedReport
		includedReport.Valid = true
	}

	habit, err := apiCfg.DB.UpdateFollowHabit(r.Context(), database.UpdateFollowHabitParams{
		HabitID: habitId,
		UserID:  user.ID,
		PromiseFrom: sql.NullTime{
			Time:  params.PromiseFrom,
			Valid: !params.PromiseFrom.IsZero(),
		},
		PromiseEnd: sql.NullTime{
			Time:  params.PromiseEnd,
			Valid: !params.PromiseEnd.IsZero(),
		},
		IncludedReport: includedReport,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update follow habit %v", err))
		return
	}

	responseWithJSON(w, 200, habit)
}

func (apiCfg *apiConfig) handlerUnfollowHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	habitIdStr := chi.URLParam(r, "habit_id")
	habitId, err := uuid.Parse(habitIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing habit id %v", err))
		return
	}
	err = apiCfg.DB.DeleteFollowHabit(r.Context(), database.DeleteFollowHabitParams{
		UserID:  user.ID,
		HabitID: habitId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't unfollow habit %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
