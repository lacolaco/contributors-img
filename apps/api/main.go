package main

import (
	"log"

	app "contrib.rocks/apps/api/internal"
)

func main() {
	log.Fatal(app.StartServer())
}
