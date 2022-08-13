package main

import (
	"log"

	app "contrib.rocks/apps/api/internal"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(0)
	godotenv.Load()
	log.Fatal(app.StartServer())
}
