package main

import (
	"log"
	"note_app/internal/app"
	"note_app/internal/configs"
)

func main() {
	application := app.NewApp()
	application.Initialize()
	log.Fatal(application.Run(configs.Config.Port))
}
