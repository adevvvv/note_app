package config

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
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
