package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CampaignDelete(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "campaign_id")
	err := h.CampaignService.Delete(id)

	if err != nil {
		return nil, 400, err
	}

	return true, 200, nil
}
