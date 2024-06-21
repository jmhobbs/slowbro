package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"slices"

	"github.com/rs/zerolog/log"
)

const separator string = "\n========================================\n"

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dumpBody := slices.Contains([]string{"application/json", "text/plain"}, r.Header.Get("Content-Type"))
		req, err := httputil.DumpRequest(r, dumpBody)
		if err != nil {
			log.Error().Err(err).Msg("Error dumping request")
		} else {
			fmt.Println(separator + string(req))
		}
		next.ServeHTTP(w, r)
	})
}
