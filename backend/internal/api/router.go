package api

import (
	"net/http"

	"backend/internal/db"
	config "backend/internal/setup"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type RouterDeps struct {
	Repo   db.Repository
	Config *config.AppConfig
}

func NewRouter(repo db.Repository, cfg *config.AppConfig) http.Handler {
	r := mux.NewRouter()
	deps := &RouterDeps{
		Repo:   repo,
		Config: cfg,
	}

	r.HandleFunc("/transactions", deps.CreateTransactionHandler).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(r)
	return handler
}
