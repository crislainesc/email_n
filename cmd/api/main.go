package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infra/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("fail to read environment variables")
	}

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
	router.Patch("/campaigns/cancel/{campaign_id}", endpoints.HandlerError(handler.CampaignCancelPatch))

	http.ListenAndServe(":3000", router)
}
