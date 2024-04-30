package db

import (
	"beats/db/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(addr string) {
	connection, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	DB = connection
	err = connection.AutoMigrate(
		&models.User{},
		&models.Beat{},
		&models.Demo{},
		&models.Genre{},
		&models.License{},
		&models.LicenseType{},
		&models.Purchase{},
		&models.Review{},
		&models.Role{},
		&models.Snippet{},
		&models.Tag{},
		&models.Verification{},
		&models.VerificationStatus{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
}
