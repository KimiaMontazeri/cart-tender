package db

import (
	"fmt"
	"time"

	"github.com/KimiaMontazeri/cart-tender/internal/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const HealthCheckPeriod = 1 * time.Second

func WithRetry(fn func(cfg config.PostgresDB) (*gorm.DB, error), cfg config.PostgresDB) *gorm.DB {
	const maxAttempts = 60
	for i := 0; i < maxAttempts; i++ {
		db, err := fn(cfg)
		if err == nil {
			return db
		}

		logrus.Errorf("Could not connect to DB. Waiting 1 second. Reason is => %s", err.Error())
		<-time.After(HealthCheckPeriod)
	}

	panic(fmt.Sprintf("Could not connect to postgres after %d attempts", maxAttempts))
}

func Connect(cfg config.PostgresDB) (*gorm.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, int(cfg.ConnectTimeout.Seconds()),
	)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
