package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/nakatanakatana/iap-user-proxy"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrEnvReqired = errors.New("Environment Variables Require")
	wrapErrorInfo = "%w: %s"
)

type Config struct {
	ClientID     string
	ClientSecret string
	BackendURL   string
	Port         string
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Println(wrapErrorInfo+"\n", ErrEnvReqired, key)
		os.Exit(1)
	}

	return value
}

func createConfig() Config {
	clientID := getRequiredEnv("IAP_PROXY_CLIENT_ID")
	clientSecret := getRequiredEnv("IAP_PROXY_CLIENT_SECRET")
	backendURL := getRequiredEnv("IAP_PROXY_BACKEND_URL")

	port := os.Getenv("IAP_PROXY_PORT")
	if port == "" {
		port = "18000"
	}

	return Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BackendURL:   backendURL,
		Port:         port,
	}
}

func main() {
	cfg := createConfig()
	config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("http://localhost:%s/__/redirect", cfg.Port),
		Scopes:       []string{"openid", "email"},
	}

	log.Println(config.AuthCodeURL(iap.AuthCodeKey))

	backend, err := url.Parse(cfg.BackendURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	rp := iap.CreateReverseProxy(backend)
	mux := http.NewServeMux()
	mux.Handle("/", rp)
	mux.Handle("/__/login", iap.LoginHandler(config))
	mux.Handle("/__/redirect", iap.ReceiveRedirectHandler(config, rp))
	log.Fatal(http.ListenAndServe(":18000", mux))
}
