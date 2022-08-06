package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"contrib.rocks/libs/goutils/env"
	"github.com/gobuffalo/envy"
)

func main() {
	envy.Load()
	appEnv := env.FromString(envy.Get("APP_ENV", "development"))
	port := envy.Get("PORT", "8080")
	fmt.Printf("Environment: %s\n", appEnv)

	http.HandleFunc("/update-featured-repositories", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repositories, err := QueryFeaturedRepositories(ctx, appEnv)
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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
