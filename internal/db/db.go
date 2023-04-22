package db

import (
	"fmt"
	"github.com/k4zb3k/project/config"
	"github.com/k4zb3k/project/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection(cfg config.DatabaseConnConfig) (*gorm.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Dbname, cfg.Password, cfg.Sslmode)

	db, err := gorm.Open(postgres.Open(conn))
	if err != nil {
		logger.Error.Printf("%s GoPostgresConnection -> Open error", err.Error())
		return nil, err
	}

	return db, nil
}
