package transport

import (
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
		seat.POST("", h.BuySeatTransport)
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

func (h *SeatHandler) BuySeatTransport(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if seat.User == "" || seat.Movie == "" || seat.Time == "" || seat.Seat == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fulfill all data"})
		return
	}

	status, err := h.service.BuySeatSvc(&seat)
	if err != nil {
		// For DB or unknown errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case repository.StatusSuccessBought:
		c.JSON(http.StatusCreated, gin.H{"message": "Seat is bought successfully", "seat": seat})
	case repository.StatusSuccessReserved:
		c.JSON(http.StatusCreated, gin.H{"message": "Seat is booked successfully", "seat": seat})
	case repository.StatusSeatSold:
		c.JSON(http.StatusConflict, gin.H{"error": "Seat is already sold"})
	case repository.StatusReservedByAnother:
		c.JSON(http.StatusForbidden, gin.H{"error": "Seat is reserved by another user"})
	case repository.StatusSeatNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
	case repository.StatusUnknownStatus:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown seat status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}
}

func (h *SeatHandler) CreateMovie(c *gin.Context) {
	var seat repository.Hall1 //repository.Hall1
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
