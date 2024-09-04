package postgres

import (
	"caa-test/internal/config"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORM(c config.Database) *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.DataSourceName()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal().Msgf("failed to opening db conn: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Msgf("failed to get db object: %s", err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}