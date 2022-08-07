package main

import (
	"log"

	app "contrib.rocks/apps/worker/internal"
)

func main() {
	log.Fatal(app.StartServer())
}
