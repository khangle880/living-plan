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

func (apiCfg *apiConfig) handlerCreateDiary(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title       string     `json:"title,omitempty"`
		Description string     `json:"description,omitempty"`
		Icon        string     `json:"icon,omitempty"`
		Background  string     `json:"background,omitempty"`
		Images      []string   `json:"images"`
		Urls        []string   `json:"urls"`
		DatetimeExc *time.Time `json:"datetime_exc,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Params decoding error: %v", err))
		return
	}
	if params.DatetimeExc == nil {
		respondWithError(w, 400, "datetime_exc cannot nil")
		return
	}

	diary, err := apiCfg.DB.CreateDiary(r.Context(), database.CreateDiaryParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       params.Title,
		Description: sql.NullString{String: params.Description, Valid: true},
		Background:  sql.NullString{String: params.Background, Valid: true},
		Icon:        sql.NullString{String: params.Icon, Valid: true},
		Images:      params.Images,
		Urls:        params.Urls,
		DatetimeExc: *params.DatetimeExc,
		UserID:      user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create diary %v", err))
		return
	}

	responseWithJSON(w, 201, diary)
}
func (apiCfg *apiConfig) handlerUpdateDiary(w http.ResponseWriter, r *http.Request, user database.User) {
	diaryIdStr := chi.URLParam(r, "diary_id")
	diaryId, err := uuid.Parse(diaryIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing diary id %v", err))
		return
	}

	type parameters struct {
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty"`
		Icon        string    `json:"icon,omitempty"`
		Background  string    `json:"background,omitempty"`
		Images      []string  `json:"images"`
		Urls        []string  `json:"urls"`
		DatetimeExc time.Time `json:"datetime_exc,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Params decoding error: %v", err))
		return
	}

	diary, err := apiCfg.DB.UpdateDiary(r.Context(), database.UpdateDiaryParams{
		ID: diaryId,
		Title: sql.NullString{
			String: params.Title,
			Valid:  params.Title != "",
		},
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},
		Icon: sql.NullString{
			String: params.Icon,
			Valid:  params.Icon != "",
		},
		Background: sql.NullString{
			String: params.Background,
			Valid:  params.Background != "",
		},
		Images: params.Images,
		Urls:   params.Urls,
		DatetimeExc: sql.NullTime{
			Time:  params.DatetimeExc,
			Valid: !params.DatetimeExc.IsZero(),
		},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update diary %v", err))
		return
	}
	responseWithJSON(w, 200, diary)
}

func (apiCfg *apiConfig) handlerGetDiary(w http.ResponseWriter, r *http.Request, user database.User) {
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
	diaries, err := apiCfg.DB.GetDiariesByUser(r.Context(), database.GetDiariesByUserParams{
		UserID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get diary %v", err))
		return
	}
	responseWithJSON(w, 200, diaries)
}

func (apiCfg *apiConfig) handlerDeleteDiary(w http.ResponseWriter, r *http.Request, user database.User) {
	diaryIdStr := chi.URLParam(r, "diary_id")
	diaryId, err := uuid.Parse(diaryIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing diary id %v", err))
		return
	}
	err = apiCfg.DB.DeleteDiary(r.Context(), database.DeleteDiaryParams{
		UserID: user.ID,
		ID:     diaryId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete diary %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
