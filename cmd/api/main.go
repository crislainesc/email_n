package main

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/errors"
	"emailn/internal/infra/database"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	service := campaign.Service{
		Repository: &database.CampaignRepository{},
	}

	router.Post("/campaigns", func(w http.ResponseWriter, r *http.Request) {
		request := contract.NewCampaign{}
		render.DecodeJSON(r.Body, &request)
		id, err := service.Create(request)

		if err != nil {

			if errors.Is(err, internalerrors.ErrorInternal) {
				render.Status(r, 500)
			} else {
				render.Status(r, 400)
			}
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		render.Status(r, 202)
		render.JSON(w, r, map[string]string{"id": id})
	})

	http.ListenAndServe(":3000", router)
}
