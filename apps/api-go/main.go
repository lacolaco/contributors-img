package main

import (
	"log"

	app "contrib.rocks/apps/api-go/internal"
)

func main() {
	log.Fatal(app.StartServer())
}
