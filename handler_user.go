package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/khangle880/living-plan/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username" validate:"required_without=Email"`
		Fullname string `json:"fullname"`
		Password string `json:"password" validate:"required"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email" validate:"required_without=Username"`
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

	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't hash password %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Username:       params.Username,
		Fullname:       params.Fullname,
		HashedPassword: hashedPassword,
		Bio:            sql.NullString{String: params.Bio, Valid: true},
		Avatar:         sql.NullString{String: params.Avatar, Valid: true},
		Email:          params.Email,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user %v", err))
		return
	}

	responseWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Username string `json:"username"`
		Fullname string `json:"fullname"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	newUser, err := apiCfg.DB.UpdateProfile(r.Context(), database.UpdateProfileParams{
		ID:       user.ID,
		Username: sql.NullString{String: params.Username, Valid: params.Username != ""},
		Fullname: sql.NullString{String: params.Fullname, Valid: params.Fullname != ""},
		Bio:      sql.NullString{String: params.Bio, Valid: params.Bio != ""},
		Avatar:   sql.NullString{String: params.Avatar, Valid: params.Avatar != ""},
		Email:    sql.NullString{String: params.Email, Valid: params.Email != ""},
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update user %v", err))
		return
	}

	responseWithJSON(w, 200, newUser)
}

func (apiCfg *apiConfig) handlerChangePassword(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't hash password %v", err))
		return
	}

	newUser, err := apiCfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update password %v", err))
		return
	}

	responseWithJSON(w, 200, newUser)
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, 200, user)
}

func (apiConfig *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username" validate:"required_without=Email"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required_without=Username"`
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

	user, err := apiConfig.DB.GetUser(r.Context(), database.GetUserParams{
		Username: params.Username,
		Email:    params.Email,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid username or password %v", err))
		return
	}

	// compare password
	err = comparePasswords(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid username or password %v", err))
		return
	}

	responseWithJSON(w, 200, user.ApiKey)

}

func (apiCfg *apiConfig) handlerGetHabitsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
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
	habits, err := apiCfg.DB.GetHabitsByUser(r.Context(), database.GetHabitsByUserParams{
		UserID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get habits %v", err))
		return
	}
	responseWithJSON(w, 200, habits)
}

func (apiCfg *apiConfig) handlerGetProjectsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
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
	projects, err := apiCfg.DB.GetProjectsByUser(r.Context(), database.GetProjectsByUserParams{
		UserID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get projects by user %v", err))
		return
	}
	responseWithJSON(w, 200, projects)
}

func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	err := apiCfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete user %v", err))
		return
	}
	responseWithJSON(w, 200, struct{}{})
}
