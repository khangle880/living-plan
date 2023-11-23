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

func (apiCfg *apiConfig) handlerCreateProject(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title       string `json:"title"`
		Purpose     string `json:"purpose"`
		Description string `json:"description"`
		Background  string `json:"background"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	project, err := apiCfg.DB.CreateProject(r.Context(), database.CreateProjectParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title:     params.Title,
		Creator:   user.ID,
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},

		Purpose: sql.NullString{
			String: params.Purpose,
			Valid:  params.Purpose != "",
		},
		Background: sql.NullString{
			String: params.Background,
			Valid:  params.Background != "",
		},
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user %v", err))
		return
	}

	_, err = apiCfg.DB.CreateFollowProject(r.Context(), database.CreateFollowProjectParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		ProjectID: project.ID,
	})

	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Couldn't follow project %v", err))
		return
	}

	responseWithJSON(w, 201, project)
}

func (apiCfg *apiConfig) handlerUpdateProject(w http.ResponseWriter, r *http.Request, user database.User) {
	projectIdStr := chi.URLParam(r, "project_id")
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing project id %v", err))
		return
	}

	type parameters struct {
		Title       string `json:"title"`
		Purpose     string `json:"purpose"`
		Description string `json:"description"`
		Background  string `json:"background"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	newProject, err := apiCfg.DB.UpdateProject(r.Context(),
		database.UpdateProjectParams{
			ID:      projectId,
			Creator: user.ID,
			Title: sql.NullString{
				String: params.Title,
				Valid:  params.Title != "",
			},
			Description: sql.NullString{
				String: params.Description,
				Valid:  params.Description != "",
			},

			Purpose: sql.NullString{
				String: params.Purpose,
				Valid:  params.Purpose != "",
			},
			Background: sql.NullString{
				String: params.Background,
				Valid:  params.Background != "",
			},
		})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update project %v", err))
		return
	}
	responseWithJSON(w, 200, newProject)
}

func (apiCfg *apiConfig) handlerGetProjects(w http.ResponseWriter, r *http.Request) {
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
	projects, err := apiCfg.DB.GetProjects(r.Context(), database.GetProjectsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get projects %v", err))
		return
	}
	responseWithJSON(w, 200, projects)
}

func (apiCfg *apiConfig) handlerDeleteProject(w http.ResponseWriter, r *http.Request, user database.User) {
	projectIdStr := chi.URLParam(r, "project_id")
	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing project id %v", err))
		return
	}
	err = apiCfg.DB.DeleteProject(r.Context(), database.DeleteProjectParams{
		ID:      projectId,
		Creator: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete project %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
