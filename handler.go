package iap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"golang.org/x/oauth2"
)

const (
	AuthCodeKey = "iapuserproxy"
)

func ReceiveRedirectHandler(cfg oauth2.Config, target *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query["state"][0] != AuthCodeKey {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		code := query["code"][0]
		if code != "" {
			ctx := context.Background()

			token, err := cfg.Exchange(ctx, code)
			if err != nil {
				log.Println("receive error:", err)
			}

			target.Transport = &Transport{
				Source: cfg.TokenSource(ctx, token),
			}
		}

		fmt.Fprintf(w, "OK")
		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(cfg oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCodeURL := cfg.AuthCodeURL(AuthCodeKey)
		log.Println(authCodeURL)
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	}
}
