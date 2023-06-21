package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infra/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("fail to read environment variables")
	}
}

func main() {
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

	router.Post("/campaigns", endpoints.HandlerError(handler.CampaignPost))
	router.Get("/campaigns/{campaign_id}", endpoints.HandlerError(handler.CampaignGetById))

	http.ListenAndServe(":3000", router)
}
