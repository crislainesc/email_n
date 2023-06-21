package endpoints

import (
	internalerrors "emailn/internal/errors"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type EndpointFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func HandlerError(endpointFunc EndpointFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj, status, err := endpointFunc(w, r)

		if err != nil {
			if errors.Is(err, internalerrors.ErrorInternal) {
				render.Status(r, 500)
			} else {
				render.Status(r, 400)
			}
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		render.Status(r, status)
		if obj != nil {
			render.JSON(w, r, obj)
		}
	})
}
