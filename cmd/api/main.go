package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infra/database"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

// loadEnv loads the environment variables from the .env file.
func loadEnv(envFile string) {
	err := godotenv.Load(dir(envFile))
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}
}

/*
dir returns the absolute path of the given environment file (envFile) in the Go module's
root directory. It searches for the 'go.mod' file from the current working directory upwards
and appends the envFile to the directory containing 'go.mod'.
It panics if it fails to find the 'go.mod' file.
*/
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

func main() {
	loadEnv(".env")
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	db := database.NewDatabase()
	campaignService := campaign.ServiceImp{
		Repository: &database.CampaignRepository{Database: db},
	}
	handler := endpoints.Handler{CampaignService: &campaignService}

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.Route("/campaigns", func(r chi.Router) {
		r.Use(endpoints.Auth)

		r.Post("/", endpoints.HandlerError(handler.CampaignPost))
		r.Get("/{campaign_id}", endpoints.HandlerError(handler.CampaignGetById))
		r.Delete("/delete/{campaign_id}", endpoints.HandlerError(handler.CampaignDelete))
	})

	http.ListenAndServe(":3000", router)
}
