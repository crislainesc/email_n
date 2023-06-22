package endpoints

import (
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/render"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.Status(r, 404)
			render.JSON(w, r, map[string]string{"error": "request does not contain an authorization header"})
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)
		client_provider := os.Getenv("KEYCLOCK_CLIENT_PROVIDER")
		client_id := os.Getenv("KEYCLOCK_CLIENT_ID")

		provider, err := oidc.NewProvider(r.Context(), "http://localhost:8080/realms/"+client_provider)

		if err != nil {
			println(err.Error())

			render.Status(r, 500)
			render.JSON(w, r, map[string]string{"error": "error to connect to server"})
			return

		}

		verifier := provider.Verifier(&oidc.Config{ClientID: client_id})
		_, err = verifier.Verify(r.Context(), token)

		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "invalid token"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
