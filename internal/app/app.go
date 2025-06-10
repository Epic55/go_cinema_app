package app

import (
	"fmt"
	"go_cinema_app/internal/logging"
	"go_cinema_app/internal/middleware"
	"go_cinema_app/internal/repository"
	"go_cinema_app/internal/services"
	"go_cinema_app/internal/transport"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	config := LoadConfig()
	logger := logging.NewLogger()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.DB.Host, config.DB.User, config.DB.Password,
		config.DB.Name, config.DB.Port, config.DB.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}

	seatRepo := repository.NewSeatRepository(db)
	seatService := services.NewSeatService(seatRepo)
	seatHandler := transport.NewSeatHandler(seatService)

	r := gin.Default()
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.AuthMiddleware(config.Auth.APIKey))

	seatHandler.RegisterRoutes(r)

	logger.Infof("Starting server on :%s", config.ServerPort)
	r.Run("localhost:" + config.ServerPort)

}
