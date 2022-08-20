package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"contrib.rocks/libs/go/env"
	"github.com/joho/godotenv"
)

func StartServer() error {
	godotenv.Load()
	appEnv := env.FromString(os.Getenv("APP_ENV"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}
	fmt.Printf("Environment: %s\n", appEnv)

	http.HandleFunc("/update-featured-repositories", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repositories, err := QueryFeaturedRepositories(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = SaveFeaturedRepositories(ctx, appEnv, repositories, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		body, _ := json.MarshalIndent(repositories, "", "  ")
		w.Write(body)
	})

	fmt.Printf("Listening on http://localhost:%s\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
