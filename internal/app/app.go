package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"note_app/internal/configs"
	"note_app/internal/handlers"
	"note_app/internal/repository"
	"note_app/internal/services"
	"os"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() {
	err := initConfig()
	if err != nil {
		log.Fatal("Ошибка чтения файла конфигурации:", err)
	}

	dbConfig := configs.DBConfig{
		Host:     configs.Config.DB.Host,
		Port:     configs.Config.DB.Port,
		User:     configs.Config.DB.User,
		Password: configs.Config.DB.Password,
	}

	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	http.HandleFunc("/signup", handlers.SignupHandler(userService))
	http.HandleFunc("/login", handlers.LoginHandler(userService, configs.Config.JWTSecret))

}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func initConfig() error {
	data, err := os.ReadFile("internal/configs/config.yaml")
	if err != nil {
		return err
	}

	var conf configs.Configuration
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return err
	}

	configs.Config = conf
	return nil
}
