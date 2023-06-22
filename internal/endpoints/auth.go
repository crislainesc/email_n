package endpoints

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
)

type contextKey string

const EmailKey contextKey = "email"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.Status(r, 404)
			render.JSON(w, r, map[string]string{"error": "request does not contain an authorization header"})
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
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
		_, err = verifier.Verify(r.Context(), tokenString)

		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "invalid token"})
			return
		}

		token, _ := jwt.Parse(tokenString, nil)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		ctx := context.WithValue(r.Context(), EmailKey, email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
