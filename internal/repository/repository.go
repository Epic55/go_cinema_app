package repository

import (
	"cinema_app_gpt/internal/services"
	"log"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	logger *slog.Logger
	pool   *gorm.DB
}

func NewRepository(logger *slog.Logger, pool *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}

func InitDBConn() *gorm.DB {

	var db *gorm.DB
	dsn := "host=localhost user=postgres password=1 dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	db.AutoMigrate(&services.User{}, &services.Movie{}, &services.Showtime{}, &services.Seat{})

	return db
}
