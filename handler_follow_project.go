package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/khangle880/living-plan/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFollowProject(w http.ResponseWriter, r *http.Request, user database.User) {
	projectIdStr := chi.URLParam(r, "project_id")
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing project id %v", err))
		return
	}

	followProject, err := apiCfg.DB.CreateFollowProject(r.Context(), database.CreateFollowProjectParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ProjectID: projectId,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Couldn't follow project %v", err))
		return
	}

	responseWithJSON(w, 201, followProject)
}

func (apiCfg *apiConfig) handlerUnfollowProject(w http.ResponseWriter, r *http.Request, user database.User) {
	projectIdStr := chi.URLParam(r, "project_id")
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing project id %v", err))
		return
	}
	err = apiCfg.DB.DeleteFollowProject(r.Context(), database.DeleteFollowProjectParams{
		UserID:    user.ID,
		ProjectID: projectId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't unfollow project %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
