// main.go
package main

import (
	"cinema_app_gpt/internal/app"
	"cinema_app_gpt/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	ap, err := app.NewApp()
	if err != nil {
		//slog.Error("failed to create app", "error", err)
		return
	}

	if err = ap.Run(); err != nil {
		//slog.Error("failed to run app", "error", err)
	}

	r := gin.Default()

	r.POST("/users", services.CreateUser)
	r.POST("/movies", services.CreateMovie)
	r.POST("/showtimes", services.CreateShowtime)
	r.POST("/book", services.BookSeat)

	r.Run("localhost:8080")
}
