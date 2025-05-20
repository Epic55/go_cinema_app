package transport

import (
	"fmt"
	"microservice_template/internal/repository"
	"microservice_template/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	service services.SeatService
}

func NewSeatHandler(s services.SeatService) *SeatHandler {
	return &SeatHandler{service: s}
}

func (h *SeatHandler) RegisterRoutes(r *gin.Engine) {
	seat := r.Group("/seat")
	{
		seat.GET("", h.GetSeats)
		seat.POST("", h.BuySeat)
	}

	movie := r.Group("/createmovie")
	{
		movie.POST("", h.CreateMovie)
	}
}

func (h *SeatHandler) GetSeats(c *gin.Context) {
	seat, err := h.service.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve seat"})
		return
	}
	c.JSON(http.StatusOK, seat)
}

func (h *SeatHandler) BuySeat(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if seat.User == "" || seat.Movie == "" || seat.Time == "" || seat.Seat == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fulfill all data"})
		return
	}

	fmt.Println("------BuySeat seat - ", seat)
	if err := h.service.BuySeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not buy seat"})
		return
	}
	c.JSON(http.StatusCreated, seat)
}

func (h *SeatHandler) CreateMovie(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.CreateMovie(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create movie"})
		return
	}

	c.JSON(http.StatusCreated, seat)
}
