package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/khangle880/living-plan/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

var validate *validator.Validate

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the enviroment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	v1Router := chi.NewRouter()
	// Users
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/register", apiCfg.handlerCreateUser)
	v1Router.Get("/login", apiCfg.handlerLogin)
	v1Router.Put("/me/profile", apiCfg.middlewareAuth(apiCfg.handlerUpdateUser))
	v1Router.Put("/change_password", apiCfg.middlewareAuth(apiCfg.handlerChangePassword))
	v1Router.Delete("/users", apiCfg.middlewareAuth(apiCfg.handlerDeleteUser))
	v1Router.Get("/me/projects", apiCfg.middlewareAuth(apiCfg.handlerGetProjectsByUser))
	v1Router.Get("/me/habits", apiCfg.middlewareAuth(apiCfg.handlerGetHabitsByUser))

	// Tasks
	v1Router.Post("/tasks", apiCfg.middlewareAuth(apiCfg.handlerCreateTask))
	v1Router.Get("/tasks", apiCfg.middlewareAuth(apiCfg.handlerGetTasks))
	v1Router.Put("/tasks/{task_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateTask))
	v1Router.Put("/tasks/{task_id}/change_status", apiCfg.middlewareAuth(apiCfg.handlerChangeStatus))
	v1Router.Delete("/tasks/{task_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteTask))

	// projects
	v1Router.Post("/projects", apiCfg.middlewareAuth(apiCfg.handlerCreateProject))
	v1Router.Get("/projects", apiCfg.handlerGetProjects)
	v1Router.Put("/projects/{project_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateProject))
	v1Router.Delete("/projects/{project_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteProject))
	v1Router.Post("/projects/follow/{project_id}", apiCfg.middlewareAuth(apiCfg.handlerCreateFollowProject))
	v1Router.Delete("/projects/unfollow/{project_id}", apiCfg.middlewareAuth(apiCfg.handlerUnfollowProject))

	// habits
	v1Router.Post("/habits", apiCfg.middlewareAuth(apiCfg.handlerCreateHabit))
	v1Router.Get("/habits", apiCfg.handlerGetHabits)
	v1Router.Put("/habits/{habit_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateHabit))
	v1Router.Delete("/habits/{habit_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteHabit))
	v1Router.Post("/habits/follow/{habit_id}", apiCfg.middlewareAuth(apiCfg.handlerCreateFollowHabit))
	v1Router.Put("/habits/follow/{habit_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateFollowHabit))
	v1Router.Delete("/habits/unfollow/{habit_id}", apiCfg.middlewareAuth(apiCfg.handlerUnfollowHabit))

	// diaries
	v1Router.Post("/diaries", apiCfg.middlewareAuth(apiCfg.handlerCreateDiary))
	v1Router.Get("/diaries", apiCfg.middlewareAuth(apiCfg.handlerGetDiary))
	v1Router.Put("/diaries/{diary_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateDiary))
	v1Router.Delete("/diaries/{diary_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteDiary))

	router.Mount("/v1", v1Router)

	fmt.Printf("Server starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
