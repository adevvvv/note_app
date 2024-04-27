package configs

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Configuration struct {
	Port      string   `yaml:"port"`
	JWTSecret string   `yaml:"jwtSecret"`
	DB        DBConfig `yaml:"db"`
}

func (c *DBConfig) Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", c.Host, c.Port, c.User, c.Password)
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

var Config Configuration
