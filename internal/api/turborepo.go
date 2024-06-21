package api

import (
	"net/http"
	"net/url"
)

func Login(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirect := r.URL.Query().Get("redirect_uri")

		u, err := url.Parse(redirect)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		q := u.Query()
		q.Set("token", token)
		q.Set("name", "Slowbro")
		u.RawQuery = q.Encode()

		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	}
}

func LoginSuccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success! You may close this window and return to the CLI."))
}
