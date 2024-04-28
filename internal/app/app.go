package app

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"note_app/internal/config"
	"note_app/internal/handlers"
	"note_app/internal/repository"
	"note_app/internal/services"
	"os"
)

// App представляет собой приложение, которое содержит маршрутизатор Gin.
type App struct {
	Router *gin.Engine
}

// NewApp создает новый экземпляр приложения.
func NewApp() *App {
	return &App{}
}

// Initialize инициализирует приложение, подготавливает маршрутизатор и подключается к базе данных.
func (a *App) Initialize() error {
	err := initConfig()
	if err != nil {
		return err
	}

	db, err := config.Config.DB.Connect()
	if err != nil {
		return err
	}

	userService := services.NewUserService(repository.NewUserRepository(db))
	noteService := services.NewNoteService(repository.NewNoteRepository(db))

	// Инициализируем маршрутизатор Gin
	a.Router = gin.Default()

	// Переключаемся в режим выпуска в производственной среде
	gin.SetMode(gin.ReleaseMode)

	// Используем обработчики Gin
	a.initHandlers(userService, &noteService)

	return nil
}

// initHandlers инициализирует обработчики запросов и добавляет их к маршрутизатору.
func (a *App) initHandlers(userService *services.UserService, noteService *services.NoteService) {
	signUpHandler := handlers.NewSignupHandler(userService)
	loginHandler := handlers.NewLoginHandler(userService, config.Config.JWTSecret)
	noteHandler := handlers.NewNoteHandler(*noteService, userService, config.Config.JWTSecret)
	editNoteHandler := handlers.EditNote(*noteService, userService, config.Config.JWTSecret)
	deleteNoteHandler := handlers.DeleteNote(*noteService, config.Config.JWTSecret)

	a.Router.POST("/signup", signUpHandler.SignUp)
	a.Router.POST("/login", loginHandler.Login)
	a.Router.POST("/note", noteHandler.AddNote)
	a.Router.PUT("/note/:id", editNoteHandler)
	a.Router.DELETE("/note/:id", deleteNoteHandler)
}

// Run запускает сервер на указанном адресе.
func (a *App) Run(addr string) error {
	return a.Router.Run(addr)
}

// initConfig инициализирует конфигурацию приложения из файла YAML.
func initConfig() error {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return err
	}

	var conf config.Configuration
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return err
	}

	config.Config = conf
	return nil
}
