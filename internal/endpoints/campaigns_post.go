package endpoints

import (
	"emailn/internal/contract"
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handler) CampaignPost(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	request := contract.NewCampaignInput{}
	render.DecodeJSON(r.Body, &request)
	id, err := h.CampaignService.Create(request)

	return map[string]string{"id": id}, http.StatusCreated, err
}
