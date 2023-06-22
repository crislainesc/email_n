package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CampaignCancelPatch(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "campaign_id")

	err := h.CampaignService.Cancel(id)

	if err != nil {
		return nil, 400, err
	}

	return true, 200, nil
}
