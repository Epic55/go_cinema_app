package transport

import (
	"cinema_app_gpt/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	logger *slog.Logger
	DB     *gorm.DB
	//repo   Repository
}

func NewHandlers(logger *slog.Logger, db *gorm.DB) *Handlers { //repo Repository
	return &Handlers{
		logger: logger,
		DB:     db,
		//repo:   repo,
	}
}

func (a *App) InitServer(db *gorm.DB) error {
	a.repository = repository.NewRepository(a.logger, a.dbConn)

	err := a.repository.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	a.handlers = NewHandlers(a.logger, db)

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/v1/bookings/:workshop_id", a.handlers.CreateBooking)
	router.GET("/api/v1/bookings/:workshop_id", a.handlers.ListBookings)

	a.router = router

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", "8080"),
		Handler: router,
	}

	a.http = server
	a.closers = append(a.closers, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return a.http.Shutdown(ctx)
	})

	return nil
}
