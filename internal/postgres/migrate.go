package postgres

import (
	"caa-test/internal/entity"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.AgentAllocationQueue{})
	if err != nil {
		log.Fatal().Msgf("failed to run migration: %s", err.Error())
	}
}