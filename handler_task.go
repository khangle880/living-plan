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

func (apiCfg *apiConfig) handlerCreateTask(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title       string              `json:"title" validate:"required"`
		DatetimeExc *time.Time          `json:"datetime_exc" validate:"required"`
		Purpose     string              `json:"purpose"`
		Description string              `json:"description"`
		Images      []string            `json:"images"`
		Urls        []string            `json:"urls"`
		Status      database.StatusEnum `json:"status" validate:"required"`
		ProjectId   uuid.NullUUID       `json:"project_id,omitempty"`
		ParentId    uuid.NullUUID       `json:"parent_id,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	err = validate.Struct(params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Validation error: %v", err))
		return
	}
	if params.DatetimeExc == nil {
		respondWithError(w, 400, "datetime_exc can not nil")
		return
	}

	task, err := apiCfg.DB.CreateTask(r.Context(), database.CreateTaskParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		DatetimeExc: *params.DatetimeExc,
		Title:       params.Title,
		Purpose:     sql.NullString{String: params.Purpose, Valid: params.Purpose != ""},
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		Images:      params.Images,
		Urls:        params.Urls,
		Status:      params.Status,
		UserID:      user.ID,
		ProjectID:   params.ProjectId,
		ParentID:    params.ParentId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create task %v", err))
		return
	}

	responseWithJSON(w, 201, task)
}

func (apiCfg *apiConfig) handlerUpdateTask(w http.ResponseWriter, r *http.Request, user database.User) {
	taskIdStr := chi.URLParam(r, "task_id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing task id %v", err))
		return
	}

	type parameters struct {
		Title       string              `json:"title"`
		DatetimeExc time.Time           `json:"datetime_exc"`
		Purpose     string              `json:"purpose"`
		Description string              `json:"description"`
		Images      []string            `json:"images"`
		Urls        []string            `json:"urls"`
		Status      database.StatusEnum `json:"status"`
		ProjectId   uuid.NullUUID       `json:"project_id"`
		ParentId    uuid.NullUUID       `json:"parent_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	task, err := apiCfg.DB.UpdateTask(r.Context(), database.UpdateTaskParams{
		ID: taskId,
		DatetimeExc: sql.NullTime{
			Time:  params.DatetimeExc,
			Valid: !params.DatetimeExc.IsZero(),
		},
		Title:       sql.NullString{String: params.Title, Valid: params.Title != ""},
		Purpose:     sql.NullString{String: params.Purpose, Valid: params.Purpose != ""},
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		Images:      params.Images,
		Urls:        params.Urls,
		ProjectID:   params.ProjectId,
		ParentID:    params.ParentId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update task %v", err))
		return
	}

	responseWithJSON(w, 200, task)
}

func (apiCfg *apiConfig) handlerChangeStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	taskIdStr := chi.URLParam(r, "task_id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing task id %v", err))
		return
	}

	type parameters struct {
		Status database.StatusEnum `json:"status"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	tasks, err := apiCfg.DB.UpdateStatus(r.Context(), database.UpdateStatusParams{
		ID:     taskId,
		Status: params.Status,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't change task status %v", err))
		return
	}
	responseWithJSON(w, 200, tasks)
}

func (apiCfg *apiConfig) handlerGetTasks(w http.ResponseWriter, r *http.Request, user database.User) {
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

	tasks, err := apiCfg.DB.GetTasksByUser(r.Context(), database.GetTasksByUserParams{
		UserID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get tasks %v", err))
		return
	}
	responseWithJSON(w, 200, tasks)
}

func (apiCfg *apiConfig) handlerDeleteTask(w http.ResponseWriter, r *http.Request, user database.User) {
	taskIdStr := chi.URLParam(r, "task_id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing task id %v", err))
		return
	}
	err = apiCfg.DB.DeleteTask(r.Context(), database.DeleteTaskParams{
		UserID: user.ID,
		ID:     taskId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete task %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
