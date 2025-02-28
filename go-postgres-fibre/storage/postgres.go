package storage

import (
	// "fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config string) (*gorm.DB, error) {
	// dsn := fmt.Sprintf("host=%s port=%s password=%s user=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Password, config.User, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
