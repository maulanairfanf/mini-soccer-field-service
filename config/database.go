package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	config := Config

	logrus.Infof("Connecting to DB with user: %s, host: %s, database: %s, password: %s, port: %s", config.Database.Username, config.Database.Host, config.Database.Name, config.Database.Password, config.Database.Port)

	encodedPassword := url.QueryEscape(config.Database.Password)

	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.Username,
		encodedPassword,
		config.Database.Host, // use localhost for local development | use postgres for docker
		config.Database.Port, // use 5433 for local development | use 5432 for docker
		config.Database.Name,
	)

	logrus.Infof("uri %v", uri)
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		logrus.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("Failed to get database connection: %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)
	return db, nil
}
