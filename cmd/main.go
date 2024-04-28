package main

import (
	"log"
	"note_app/internal/app"
	"note_app/internal/config"
)

func main() {
	application := app.NewApp()
	application.Initialize()
	log.Fatal(application.Run(config.Config.Port))
}
