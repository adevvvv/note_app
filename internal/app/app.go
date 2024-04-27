package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"note_app/internal/config"
	"note_app/internal/handlers"
	"note_app/internal/repository"
	"note_app/internal/services"
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

	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/signup", handlers.SignupHandler(userService)).Methods("POST")
}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.Router)
}
