package services

import (
	"cinema_app_gpt/internal/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Repository interface {
// 	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
// 	ListBookings(ctx context.Context, workshopID int64) ([]*domain.Booking, error)
// }

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repository.Db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func CreateMovie(c *gin.Context) {
	var movie Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&movie)
	c.JSON(http.StatusCreated, movie)
}

func CreateShowtime(c *gin.Context) {
	var req ShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	show := Showtime{MovieID: req.MovieID, Start: req.Start}
	db.Create(&show)

	for i := 1; i <= 10; i++ {
		seat := Seat{
			SeatNumber: fmt.Sprintf("A%d", i),
			IsBooked:   false,
			ShowtimeID: show.ID,
		}
		db.Create(&seat)
	}

	c.JSON(http.StatusCreated, show)
}

func BookSeat(c *gin.Context) {
	var req BookSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var seat Seat
	result := db.Where("showtime_id = ? AND seat_number = ? AND is_booked = false", req.ShowtimeID, req.SeatNumber).First(&seat)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seat not available"})
		return
	}

	seat.IsBooked = true
	seat.UserID = &req.UserID
	db.Save(&seat)

	c.JSON(http.StatusOK, gin.H{"message": "Seat booked", "seat": seat})
}
