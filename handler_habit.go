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

func (apiCfg *apiConfig) handlerCreateHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title             string    `json:"title"`
		Purpose           string    `json:"purpose"`
		Description       string    `json:"description"`
		Icon              string    `json:"icon"`
		Background        string    `json:"background"`
		Images            []string  `json:"images"`
		Urls              []string  `json:"urls"`
		TimeInDay         time.Time `json:"time_in_day"`
		LoopWeek          []int32   `json:"loop_week"`
		LoopMonth         []int32   `json:"loop_month"`
		RecommendDuration int32     `json:"recommend_duration"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Params decoding error: %v", err))
		return
	}

	habit, err := apiCfg.DB.CreateHabit(r.Context(), database.CreateHabitParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		Title:             params.Title,
		Creator:           user.ID,
		Purpose:           sql.NullString{String: params.Purpose, Valid: params.Purpose != ""},
		Description:       sql.NullString{String: params.Description, Valid: params.Description != ""},
		Background:        sql.NullString{String: params.Background, Valid: params.Background != ""},
		Icon:              sql.NullString{String: params.Icon, Valid: params.Icon != ""},
		Images:            params.Images,
		Urls:              params.Urls,
		TimeInDay:         params.TimeInDay,
		LoopWeek:          params.LoopWeek,
		LoopMonth:         params.LoopMonth,
		RecommendDuration: params.RecommendDuration,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create habit %v", err))
		return
	}

	responseWithJSON(w, 201, habit)
}

func (apiCfg *apiConfig) handlerGetHabits(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	habits, err := apiCfg.DB.GetHabits(r.Context(), database.GetHabitsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get habits %v", err))
		return
	}
	responseWithJSON(w, 200, habits)
}

func (apiCfg *apiConfig) handlerUpdateHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	habitIdStr := chi.URLParam(r, "habit_id")
	habitId, err := uuid.Parse(habitIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing habit id %v", err))
		return
	}

	type parameters struct {
		Title             string    `json:"title"`
		Purpose           string    `json:"purpose"`
		Description       string    `json:"description"`
		Icon              string    `json:"icon"`
		Background        string    `json:"background"`
		Images            []string  `json:"images"`
		Urls              []string  `json:"urls"`
		TimeInDay         time.Time `json:"time_in_day"`
		LoopWeek          []int32   `json:"loop_week"`
		LoopMonth         []int32   `json:"loop_month"`
		RecommendDuration int32     `json:"recommend_duration"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	habit, err := apiCfg.DB.UpdateHabit(r.Context(), database.UpdateHabitParams{
		ID:      habitId,
		Creator: user.ID,
		TimeInDay: sql.NullTime{
			Time:  params.TimeInDay,
			Valid: !params.TimeInDay.IsZero(),
		},
		Title:       sql.NullString{String: params.Title, Valid: params.Title != ""},
		Purpose:     sql.NullString{String: params.Purpose, Valid: params.Purpose != ""},
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		Background:  sql.NullString{String: params.Background, Valid: params.Background != ""},
		Icon:        sql.NullString{String: params.Icon, Valid: params.Icon != ""},
		Images:      params.Images,
		Urls:        params.Urls,
		LoopWeek:    params.LoopWeek,
		LoopMonth:   params.LoopMonth,
		RecommendDuration: sql.NullInt32{
			Int32: params.RecommendDuration,
			Valid: params.RecommendDuration != 0,
		},
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update habit %v", err))
		return
	}

	responseWithJSON(w, 200, habit)
}

func (apiCfg *apiConfig) handlerDeleteHabit(w http.ResponseWriter, r *http.Request, user database.User) {
	habitIdStr := chi.URLParam(r, "habit_id")
	habitId, err := uuid.Parse(habitIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing habit id %v", err))
		return
	}
	err = apiCfg.DB.DeleteHabit(r.Context(), database.DeleteHabitParams{
		ID:      habitId,
		Creator: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete habit %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
