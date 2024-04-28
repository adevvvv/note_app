package app

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/http-swagger"
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

	// Инициализируем Swagger
	a.initSwagger()

	return nil
}

// initSwagger инициализирует Swagger.
func (a *App) initSwagger() {
	// Подключаем Swagger UI
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// Добавьте инициализацию нового обработчика в метод initHandlers
func (a *App) initHandlers(userService *services.UserService, noteService *services.NoteService) {
	signUpHandler := handlers.NewSignupHandler(userService).SignUp
	signInHandler := handlers.NewSignInHandler(userService, config.Config.JWTSecret).SignIn
	noteHandler := handlers.NewNoteHandler(*noteService, userService, config.Config.JWTSecret).AddNote
	editNoteHandler := handlers.EditNoteHandler(*noteService, userService, config.Config.JWTSecret)
	deleteNoteHandler := handlers.DeleteNoteHandler(*noteService, config.Config.JWTSecret)
	getNotesHandler := handlers.GetNotesHandler(*noteService, *userService, config.Config.JWTSecret)

	a.Router.POST("/signup", signUpHandler)
	a.Router.POST("/signin", signInHandler)
	a.Router.POST("/note", noteHandler)
	a.Router.PUT("/note/:id", editNoteHandler)
	a.Router.DELETE("/note/:id", deleteNoteHandler)
	a.Router.GET("/notes", getNotesHandler)

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
