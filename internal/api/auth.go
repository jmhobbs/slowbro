package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func AuthMiddleware(token string) func(http.Handler) http.Handler {
	expectedToken := fmt.Sprintf("Bearer %s", token)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != expectedToken {
				log.Warn().
					Str("authorization", r.Header.Get("Authorization")).
					Str("remote", r.RemoteAddr).
					Str("url", r.URL.String()).
					Msg("Unauthorized")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
