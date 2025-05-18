package app

import (
	"cinema_app_gpt/internal/repository"
	"cinema_app_gpt/internal/services"
	"cinema_app_gpt/internal/transport"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	logger *logrus

	dbConn     *gorm.DB
	repository *repository.Repository

	router gin.IRouter
	http   *http.Server

	closers  []func() error
	handlers *services.Handlers
	closeCh  chan os.Signal
}

func NewApp() (*App, error) {
	var err error

	app := &App{}
	// - config
	// TODO: add config
	// - logger
	app.initLogger()
	// - db
	db := repository.InitDBConn()
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}
	// - http server
	err = transport.InitServer(db)
	if err != nil {
		return nil, fmt.Errorf("failed to init http server: %w", err)
	}
	// - graceful shutdown
	err = app.initGracefulShutdown()
	if err != nil {
		return nil, fmt.Errorf("failed to init graceful shutdown: %w", err)
	}

	return app, nil
}

func (a *App) Run() error {
	go func() {
		err := a.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("http server failed", "error", err)
		}
	}()

	<-a.closeCh
	for i := len(a.closers) - 1; i >= 0; i-- {
		err := a.closers[i]()
		if err != nil {
			a.logger.Error("failed to close resource", "i", i, "error", err)
		}
	}

	return nil
}

func (a *App) initLogger() {
	logger := slog.Default()
	a.logger = logger
}

func (a *App) initGracefulShutdown() error {
	a.closeCh = make(chan os.Signal, 1)
	signal.Notify(a.closeCh, os.Interrupt)

	return nil
}
