package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CampaignDelete(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "campaign_id")
	err := h.CampaignService.Delete(id)

	return true, http.StatusOK, err
}
