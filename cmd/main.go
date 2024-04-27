package main

import (
	"log"
	"note_app/internal/app"
)

func main() {
	application := app.NewApp()
	application.Initialize()
	log.Fatal(application.Run(":8000"))
}
