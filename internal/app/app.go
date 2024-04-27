package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"note_app/internal/config"
)

type App struct {
	Router *mux.Router
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() {
	dbConfig := config.DBConfig{
		Host:     "postgres",
		Port:     5432,
		User:     "postgres",
		Password: "password",
	}

	_, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.Router)
}
