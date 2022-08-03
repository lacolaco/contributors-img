package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "staging"
	}
	fmt.Printf("Environment: %s\n", env)

	http.HandleFunc("/update-featured-repositories", func(w http.ResponseWriter, r *http.Request) {
		ctx := SetEnvironment(r.Context(), env)
		repositories, err := QueryFeaturedRepositories(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = SaveFeaturedRepositories(ctx, repositories, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		body, _ := json.MarshalIndent(repositories, "", "  ")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
