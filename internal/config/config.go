package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DBConfig представляет конфигурацию базы данных.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// Configuration представляет общую конфигурацию приложения.
type Configuration struct {
	Port      string   `yaml:"port"`
	JWTSecret string   `yaml:"jwtSecret"`
	DB        DBConfig `yaml:"db"`
}

// Connect подключается к базе данных и возвращает объект db для выполнения запросов.
func (c *DBConfig) Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Config содержит глобальную конфигурацию приложения.
var Config Configuration
